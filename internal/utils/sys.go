package utils

import "os/exec"

func CopyToClipboard(s string) error {
	cmd := exec.Command("pbcopy")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := in.Write([]byte(s)); err != nil {
		_ = in.Close()
		return err
	}
	_ = in.Close()

	return cmd.Wait()
}
