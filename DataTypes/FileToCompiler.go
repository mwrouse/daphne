package DataTypes


import (
    "os"
)

type FileToCompile struct {
    Directory   string
    Info        os.FileInfo
}
