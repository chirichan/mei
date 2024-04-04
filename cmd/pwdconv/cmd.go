package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/chirichan/mei/version"
	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
)

type XyKey struct {
	Version int   `json:"version"`
	Key     []Key `json:"key"`
}

type Key struct {
	Name      string  `json:"name"`
	Account   string  `json:"account"`
	Password  string  `json:"password"`
	Password2 string  `json:"password2"`
	Url       string  `json:"url"`
	Note      string  `json:"note"`
	Extra     []Extra `json:"extra"`
}

type Extra struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// ChromeCSV chrome, edge csv password
type ChromeCSV struct {
	Name     string `json:"name" csv:"name"`
	URL      string `json:"url" csv:"url"`
	Username string `json:"username" csv:"username"`
	Password string `json:"password" csv:"password"`
	Note     string `json:"note" csv:"note"`
}

type PwdConvCLI struct {
	Logger *slog.Logger
}

func (m *PwdConvCLI) Root(cmd *cobra.Command, args []string) {
	if v, _ := cmd.Flags().GetBool("version"); v {
		fmt.Printf("pwdgen version is %s", version.Version)
		return
	}
}

func (m *PwdConvCLI) Csv2Xykey(cmd *cobra.Command, args []string) error {
	var csvData []ChromeCSV
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

func (m *PwdConvCLI) csv2Xykey(csvData []ChromeCSV) *XyKey {
	xyKey := &XyKey{
		Version: 1,
		Key:     make([]Key, len(csvData)),
	}
	for i, v := range csvData {
		parse, err := url.Parse(v.URL)
		if err != nil {
			m.Logger.Error("parse csv name", "name", v.Name)
		}
		xyKey.Key[i] = Key{
			Name:     parse.Hostname(),
			Account:  v.Username,
			Password: v.Password,
			Url:      v.URL,
			Note:     v.Note,
			Extra:    make([]Extra, 0),
		}
	}
	return xyKey
}

func NewCLI() *cobra.Command {
	muCLI := &PwdConvCLI{Logger: slog.Default()}
	rootCmd := &cobra.Command{
		Use:   "pwdconv",
		Short: "转换密码格式工具",
		Run:   muCLI.Root,
	}
	rootCmd.Flags().BoolP("version", "v", false, "版本")
	csv2XykeyCmd := &cobra.Command{
		Use:   "csv2xykey",
		Short: "浏览器导出的 csv 格式的密码转为 xykey 格式。",
		Args:  cobra.MaximumNArgs(1),
		RunE:  muCLI.Csv2Xykey,
	}

	rootCmd.AddCommand(
		csv2XykeyCmd,
	)
	return rootCmd
}
