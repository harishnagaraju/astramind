package storage

import (
    "fmt"
    "os"

    "github.com/harishnagaraju/astramind/internal/models"
)

func ExportSession(
    session string,
    messages []models.Message,
) error {

    file := fmt.Sprintf(
        "exports/%s.txt",
        session,
    )

    f, err := os.Create(file)

    if err != nil {
        return err
    }

    defer f.Close()

    for _, msg := range messages {

        _, err := fmt.Fprintf(
            f,
            "[%s] %s\n\n",
            msg.Role,
            msg.Content,
        )

        if err != nil {
            return err
        }
    }

    return nil
}