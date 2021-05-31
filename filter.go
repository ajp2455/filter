package filter

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"sync"
)

func Filter(stream io.Reader, threads int, predicate Predicate) (io.Reader, error) {
	r, w := io.Pipe()

	go func() {
		defer w.Close()
		scanner := bufio.NewScanner(stream)

		ch := make(chan string)

		var wg sync.WaitGroup
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for line := range ch {
					ok, err := predicate(line)
					if err != nil {
						w.CloseWithError(err)
						return
					}

					if !ok {
						continue
					}

					if _, err := io.WriteString(w, line+"\n"); err != nil {
						w.CloseWithError(err)
						return
					}
				}
			}()
		}

		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)

		wg.Wait()

		w.Close()
	}()

	return r, nil
}

type Predicate func(line string) (bool, error)

func CmdPredicate(command []string, token string) Predicate {
	return func(line string) (bool, error) {
		var parsed []string
		switch token {
		case "":
			parsed = append(append([]string{}, command...), line)
		default:
			parsed = replace(command, token, line)
		}

		cmd := exec.Command(parsed[0], parsed[1:]...)
		if err := cmd.Run(); err != nil {
			var as *exec.ExitError
			if errors.As(err, &as) {
				return false, nil
			}

			return false, err
		}

		return true, nil
	}
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
