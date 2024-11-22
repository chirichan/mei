package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/chirichan/rice"
	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"

	"github.com/chirichan/mei/cmd/pwdgen/internal/entities"
	"github.com/chirichan/mei/version"
)

const Aes256Suffix = ".aes256"

type PwdGenCLI struct {
	Logger *slog.Logger
}

func (m *PwdGenCLI) Root(cmd *cobra.Command, args []string) {
	if v, _ := cmd.Flags().GetBool("version"); v {
		fmt.Printf("pwdgen version is %s", version.Version)
		return
	}
	length, _ := cmd.Flags().GetInt("length")
	level, _ := cmd.Flags().GetInt("level")
	output, _ := cmd.Flags().GetInt("output")
	s, err := fullPassword(level, length)
	if err != nil {
		log.Fatalf("full password err: %v", err)
	}
	if output == 1 {
		if err := clipboard.WriteAll(s); err != nil {
			log.Fatalf("err: %v\n", err)
		}
	} else if output == 2 {
		fmt.Println(s)
	} else {
		log.Fatalf("output param err: not support %d\n", output)
	}
}

func (m *PwdGenCLI) Csv2Xykey(cmd *cobra.Command, args []string) error {
	var csvData []entities.ChromeCSV
	if len(args) == 0 {
		csvText, err := clipboard.ReadAll()
		if err != nil {
			return err
		}
		if err := gocsv.UnmarshalString(csvText, &csvData); err != nil {
			return err
		}
	} else {
		m.Logger.Info("csv2xykey", "args", args)
		csvFile, err := os.OpenFile(args[0], os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		defer csvFile.Close()
		if err := gocsv.UnmarshalFile(csvFile, &csvData); err != nil {
			return err
		}
	}
	xykeyData := m.csv2Xykey(csvData)
	xyKeyBytes, err := json.Marshal(xykeyData)
	if err != nil {
		return err
	}
	cwd := filepath.Dir(os.Args[0])
	m.Logger.Info("save", "dir", cwd)
	return os.WriteFile(fmt.Sprintf("xykey-%d.json", time.Now().Unix()), xyKeyBytes, 0644)
}

func (m *PwdGenCLI) csv2Xykey(csvData []entities.ChromeCSV) *entities.XyKey {
	xyKey := &entities.XyKey{
		Version: 1,
		Key:     make([]entities.Key, len(csvData)),
	}
	for i, v := range csvData {
		parse, err := url.Parse(v.URL)
		if err != nil {
			m.Logger.Error("parse csv name", "name", v.Name)
		}
		xyKey.Key[i] = entities.Key{
			Name:     parse.Hostname(),
			Account:  v.Username,
			Password: v.Password,
			Url:      v.URL,
			Note:     v.Note,
			Extra:    make([]entities.Extra, 0),
		}
	}
	return xyKey
}

func (m *PwdGenCLI) SplitFile(cmd *cobra.Command, args []string) error {

	return nil
}

func (m *PwdGenCLI) EncryptFile(cmd *cobra.Command, args []string) error {

	m.Logger.Debug("encrypt cmd", "args", args)

	genKey, _ := cmd.Flags().GetBool("genkey")
	key, _ := cmd.Flags().GetString("key")
	file, _ := cmd.Flags().GetString("file")
	text, _ := cmd.Flags().GetString("text")
	begin := time.Now()

	if genKey {
		randomHexString, err := rice.RandomHexString(32)
		fmt.Printf("生成的 AES key :\n\n%s\n\n妥善保存此密钥，此后的加密和解密都需要此密钥。\n", randomHexString)
		return err
	}

	if key == "" {
		k, ok := os.LookupEnv("MEI_AES_KEY")
		if !ok {
			return fmt.Errorf("env var MEI_AES_KEY not set")
		}
		key = k
	}

	if cmd.Flags().Changed("text") && text != "" {
		encryptText, err := rice.AESGCMEncryptText(key, text)
		if err != nil {
			return fmt.Errorf("encrypt text err: %w", err)
		}
		fmt.Println(encryptText)
		return nil
	}

	if !rice.PathExists(file) {
		return fmt.Errorf("文件或文件夹不存在, file: %s", file)
	}

	if rice.PathIsDir(file) {

		absPath, _ := filepath.Abs(file)
		parentPath := filepath.Dir(absPath)
		zipFilename := filepath.Join(parentPath, filepath.Base(file)+".zip")

		if err := ZipFolder(file, zipFilename); err != nil {
			m.Logger.Error("zip folder err", "err", err)
			return err
		}

		m.Logger.Info("zip folder success", "file", file, "cost", time.Since(begin), "zip_filename", zipFilename)

		if err := rice.AESGCMEncryptFile(key, zipFilename, zipFilename+Aes256Suffix); err != nil {
			return err
		}
		if err := os.Remove(zipFilename); err != nil {
			return err
		}
		m.Logger.Info("encrypt folder success", "cost", time.Since(begin), "output", zipFilename+Aes256Suffix)

	} else {
		if err := rice.AESGCMEncryptFile(key, file, file+Aes256Suffix); err != nil {
			return err
		}
		err := os.Remove(file)
		m.Logger.Info("encrypt file success", "cost", time.Since(begin), "output", file+Aes256Suffix)
		return err
	}

	return nil
}

// ZipFolder 压缩文件夹
func ZipFolder(sourceDir, zipFile string) error {
	// 创建目标 ZIP 文件
	zipFileWriter, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zipFileWriter.Close()

	// 创建 ZIP Writer
	zipWriter := zip.NewWriter(zipFileWriter)
	defer zipWriter.Close()

	// 遍历目录并添加文件到 ZIP
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算文件相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// 忽略目录，直接处理文件
		if info.IsDir() {
			if relPath == "." {
				// 忽略根目录自身
				return nil
			}
			// 在 ZIP 中创建目录
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		// 创建文件头
		fileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		fileHeader.Name = relPath
		fileHeader.Method = zip.Deflate // 使用压缩方式

		// 创建文件写入器
		writer, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return err
		}

		// 打开原始文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// 将文件内容复制到 ZIP 中
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func UnzipFolder(src, dest string) error {
	// 打开 ZIP 文件
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// 遍历 ZIP 文件中的每个条目
	for _, file := range r.File {
		// 构造目标文件路径
		filePath := filepath.Join(dest, file.Name)

		// 检查目标路径是否安全
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("非法文件路径: %s", filePath)
		}

		// 如果是目录，则创建
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		// 如果是文件，则解压
		err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm) // 确保文件目录存在
		if err != nil {
			return err
		}

		// 打开 ZIP 文件内的文件
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 创建目标文件
		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()

		// 将内容写入目标文件
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PwdGenCLI) DecryptFile(cmd *cobra.Command, args []string) error {
	key, _ := cmd.Flags().GetString("key")
	file, _ := cmd.Flags().GetString("file")
	text, _ := cmd.Flags().GetString("text")
	begin := time.Now()

	if key == "" {
		k, ok := os.LookupEnv("MEI_AES_KEY")
		if !ok {
			return fmt.Errorf("env var MEI_AES_KEY not set")
		}
		key = k
	}

	if text != "" {
		decryptText, err := rice.AESGCMDecryptText(key, text)
		if err != nil {
			return fmt.Errorf("decrypt text err: %w", err)
		}
		fmt.Println(decryptText)
		return nil
	}

	if !rice.PathExists(file) {
		return fmt.Errorf("文件或文件夹不存在, file: %s", file)
	}

	if !strings.HasSuffix(file, Aes256Suffix) {
		return fmt.Errorf("文件不是加密文件, file: %s", file)
	}

	if rice.PathIsDir(file) {
		m.Logger.Info("这是一个文件夹，不支持解密")

	} else {
		outputFile := strings.TrimSuffix(file, Aes256Suffix)
		if err := rice.AESGCMDecryptFile(key, file, outputFile); err != nil {
			return err
		}

		if strings.HasSuffix(outputFile, ".zip") {

			if err := UnzipFolder(outputFile, strings.TrimSuffix(outputFile, ".zip")); err != nil {
				return err
			}

			if err := os.Remove(outputFile); err != nil {
				return err
			}
		}

		err := os.Remove(file)
		m.Logger.Info("解密完成", "耗时", time.Since(begin))
		return err
	}

	return nil
}

func (m *PwdGenCLI) KillProcess(cmd *cobra.Command, args []string) error {

	pNames, err := cmd.Flags().GetStringSlice("name")
	if err != nil {
		return err
	}
	for {
		var execCmd *exec.Cmd
		goos := runtime.GOOS
		for _, processName := range pNames {
			switch goos {
			case "windows":
				// taskkill.exe /IM WeChatAppEx.exe /F
				execCmd = exec.Command("taskkill", "/IM", processName+".exe", "/F")
			case "linux", "darwin":
				execCmd = exec.Command("pkill", processName)
			default:
				return fmt.Errorf("unsupported platform")
			}
			slog.Info("killing process...", "processName", processName, "os", goos)
			err := execCmd.Run()
			if err != nil {
				slog.Error("exec cmd", "err", err)
			}
		}
		time.Sleep(time.Second)
	}
}

func NewCLI() *cobra.Command {
	muCLI := &PwdGenCLI{Logger: slog.Default()}

	rootCmd := &cobra.Command{
		Use:   "pwdgen",
		Short: "生成随机密码",
		Run:   muCLI.Root,
	}
	rootCmd.Flags().BoolP("version", "v", false, "版本")
	rootCmd.Flags().IntP("length", "n", 16, "生成的密码长度, [6, 2048]")
	rootCmd.Flags().IntP("level", "l", 4, "生成的密码强度等级, 数字越大, 强度越高, [1, 4]")
	rootCmd.Flags().IntP("output", "o", 1, "输出方式, 1: 剪贴板, 2: 控制台")

	csv2XykeyCmd := &cobra.Command{
		Use:   "csv2xykey",
		Short: "浏览器导出的 csv 格式的密码转为 xykey 格式。",
		Args:  cobra.MaximumNArgs(1),
		RunE:  muCLI.Csv2Xykey,
	}

	splitFileCmd := &cobra.Command{
		Use:   "splitfile",
		Short: "把文件分割为若干小文件。",
		RunE:  muCLI.SplitFile,
	}
	splitFileCmd.Flags().IntP("size", "s", 200, "每个文件的大小, 单位: Mb")

	encryptFileCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "加密文本或文件",
		RunE:  muCLI.EncryptFile,
	}
	encryptFileCmd.Flags().BoolP("genkey", "g", false, "生成一个 AES256 密钥")
	encryptFileCmd.Flags().StringP("key", "k", "", "加密所需的密钥。如果不指定，则从环境变量 \"MEI_AES_KEY\" 中获取")
	encryptFileCmd.Flags().StringP("file", "f", "", "要加密的文件或文件夹")
	encryptFileCmd.Flags().StringP("text", "t", "", "要加密的文本")
	encryptFileCmd.Flags().StringP("ignore", "i", "", "ignore 文件【暂未实现】")

	decryptFileCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "解密文本或文件",
		RunE:  muCLI.DecryptFile,
	}
	decryptFileCmd.Flags().StringP("key", "k", "", "解密所需的密钥。如果不指定，则从环境变量 \"MEI_AES_KEY\" 中获取")
	decryptFileCmd.Flags().StringP("file", "f", "", "要解密的文件或文件夹")
	decryptFileCmd.Flags().StringP("text", "t", "", "要解密的文本")
	decryptFileCmd.MarkFlagsOneRequired("file", "text")

	killCmd := &cobra.Command{
		Use:   "kill",
		Short: "杀掉指定的进程",
		RunE:  muCLI.KillProcess,
	}
	killCmd.Flags().StringP("name", "n", "", "进程名列表")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "版本",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("pwdgen %s", version.Version)
		},
	}

	rootCmd.AddCommand(
		csv2XykeyCmd,
		splitFileCmd,
		encryptFileCmd,
		decryptFileCmd,
		killCmd,
		versionCmd,
	)
	return rootCmd
}
