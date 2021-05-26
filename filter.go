package filter

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
)

func NewFilter(stream io.Reader, command []string,
	options ...FilterOption) (io.Reader, error) {

	if len(command) == 0 {
		return nil, errors.New("command must not be empty")
	}

	var f filter
	for _, option := range options {
		option(&f)
	}

	r, w := io.Pipe()

	go func() {
		defer w.Close()

		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			line := scanner.Text()

			var parsed []string

			switch f.token {
			case "":
				parsed = append(append([]string{}, command...), line)
			default:
				parsed = replace(command, f.token, line)
			}

			cmd := exec.Command(parsed[0], parsed[1:]...)
			if err := cmd.Run(); err != nil {
				var as *exec.ExitError
				if errors.As(err, &as) {
					continue
				}

				w.CloseWithError(err)
				return
			}

			if _, err := io.WriteString(w, line+"\n"); err != nil {
				w.CloseWithError(err)
				return
			}
		}
	}()

	return r, nil
}

type Predacate func(line string) (bool, error)

type FilterOption func(f *filter)

func WithToken(token string) FilterOption {
	return func(f *filter) {
		f.token = token
	}
}

type filter struct {
	token string
}

func replace(slice []string, find, replace string) []string {
	clone := make([]string, 0, len(slice))
	for _, elem := range slice {
		switch elem {
		case find:
			clone = append(clone, replace)
		default:
			clone = append(clone, elem)
		}
	}

	return clone
}
