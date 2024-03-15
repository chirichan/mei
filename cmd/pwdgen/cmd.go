package main

import (
	"fmt"
	"log"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

type PwdgenCLI struct{}

func (m *PwdgenCLI) Root(cmd *cobra.Command, args []string) {
	if version, _ := cmd.Flags().GetBool("version"); version {
		fmt.Printf("pwdgen version is %s", Version)
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

func NewCLI() *cobra.Command {
	muCLI := &PwdgenCLI{}
	rootCmd := &cobra.Command{
		Use:   "pwdgen",
		Short: "生成随机密码",
		Run:   muCLI.Root,
	}
	rootCmd.Flags().BoolP("version", "v", false, "版本")
	rootCmd.Flags().IntP("length", "n", 16, "生成的密码长度, [6, 2048]")
	rootCmd.Flags().IntP("level", "l", 4, "生成的密码强度等级, 数字越大, 强度越高, [1, 4]")
	rootCmd.Flags().IntP("output", "o", 1, "输出方式, 1: 剪贴板, 2: 控制台")

	return rootCmd
}
