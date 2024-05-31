package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/chirichan/mei/cmd/pwdgen/internal/entities"
	"github.com/chirichan/mei/version"
	"github.com/chirichan/rice"
	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
)

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
	genKey, _ := cmd.Flags().GetBool("genkey")
	key, _ := cmd.Flags().GetString("key")
	file, _ := cmd.Flags().GetString("file")
	text, _ := cmd.Flags().GetString("text")

	if genKey {
		randomHexString, err := rice.RandomHexString(32)
		fmt.Println(randomHexString)
		return err
	}

	if key == "" {
		k, ok := os.LookupEnv("MEI_AES_KEY")
		if !ok {
			return fmt.Errorf("env var MEI_AES_KEY not set")
		}
		key = k
	}

	m.Logger.Info("encrypt", "key", key, "file", file, "text", text)

	if text != "" {
		encryptText, err := rice.AESGCMEncryptText(key, text)
		if err != nil {
			return fmt.Errorf("encrypt text err: %w", err)
		}
		fmt.Println(encryptText)
	}

	return nil
}

func (m *PwdGenCLI) DecryptFile(cmd *cobra.Command, args []string) error {
	key, _ := cmd.Flags().GetString("key")
	file, _ := cmd.Flags().GetString("file")
	text, _ := cmd.Flags().GetString("text")

	if key == "" {
		k, ok := os.LookupEnv("MEI_AES_KEY")
		if !ok {
			return fmt.Errorf("env var MEI_AES_KEY not set")
		}
		key = k
	}

	m.Logger.Info("decrypt", "key", key, "file", file, "text", text)

	if text != "" {
		decryptText, err := rice.AESGCMDecryptText(key, text)
		if err != nil {
			return fmt.Errorf("decrypt text err: %w", err)
		}
		fmt.Println(decryptText)
	}
	return nil
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
	splitFileCmd.Flags().IntP("size", "s", 200, "每个文件的大小，单位：Mb")

	encryptFileCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "加密文本或文件",
		RunE:  muCLI.EncryptFile,
	}
	encryptFileCmd.Flags().BoolP("genkey", "g", false, "生成一个 AES256 密钥")
	encryptFileCmd.Flags().StringP("key", "k", "", "加密所需的密钥。如果不指定，则从环境变量 \"MEI_AES_KEY\" 中获取")
	encryptFileCmd.Flags().StringP("file", "f", "", "要加密的文件或文件夹")
	encryptFileCmd.Flags().StringP("text", "t", "", "要加密的文本")

	decryptFileCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "解密文本或文件",
		RunE:  muCLI.DecryptFile,
	}
	decryptFileCmd.Flags().StringP("key", "k", "", "解密所需的密钥。如果不指定，则从环境变量 \"MEI_AES_KEY\" 中获取")
	decryptFileCmd.Flags().StringP("file", "f", "", "要解密的文件或文件夹")
	decryptFileCmd.Flags().StringP("text", "t", "", "要解密的文本")

	rootCmd.AddCommand(
		csv2XykeyCmd,
		splitFileCmd,
		encryptFileCmd,
		decryptFileCmd,
	)
	return rootCmd
}
