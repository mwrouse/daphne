/**
 * This package is for file system helpers
 */
package FileSystem

import (
    "daphne/Helpers"
    "daphne/Errors"
    "daphne/DataTypes"
    "os"
    "bufio"
    "io"
    "io/ioutil"
    "fmt"
)


/**
 * Name.........: ReadFile
 * Parameters...: path (string) - path to the file to read
 * Return.......: []string - array of all the lines in the file
 *                error - any errors
 * Description..: Reads the contents of a file
 */
func ReadFile(path string) ([]string, Errors.Error) {
  // Check if the file exists
  if !FileExists(path) {
    return nil, Errors.NewFatal("File ", path, " does not exist")
  }

  f, err := os.Open(path)
  if err != nil {
    return nil, Errors.NewFatal(err.Error()) // Error
  }

  defer f.Close()

  var contents []string
  reader := bufio.NewScanner(f)

  for reader.Scan() {
    contents = append(contents, reader.Text())
  }

  daphneErr := Errors.None()
  err = reader.Err()
  if err != nil {
      daphneErr = Errors.NewFatal(err.Error())
  }
  return contents, daphneErr
}


/**
 * Name.........: WriteFile
 * Parameters...: path (string) - path to the file to write to
 *                contents ([]string) - array of lines to write
 * Return.......: error - any errors
 * Description..: Writes to a file
 */
func WriteFile (path string, contents []string) (Errors.Error) {
    var f *os.File

    if Helpers.Trim(path) == "" {
        return Errors.NewFatal("No path specified for WriteFile")
    }

    path = Helpers.Replace(path, ".\\", "")
    if FileExists(path) {
        os.Remove(path)
    }

    // Get the folder name and create the folder
    splitPath := Helpers.Split(path, "\\")
    folder := Helpers.Join(splitPath[:len(splitPath) - 1], "\\")
    err1 := os.MkdirAll(folder, 0777)

    if err1 != nil {
        return Errors.NewFatal("MEow" + err1.Error())
    }

    f, err1 = os.Create(path)
    if err1 != nil {
        return Errors.NewFatal(err1.Error())
    }

    defer f.Close()

    writer := bufio.NewWriter(f)

    for _, line := range contents {
        fmt.Fprintln(writer, line)
    }

    writer.Flush()

    return Errors.None()
}


// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (Errors.Error) {
    sfi, err := os.Stat(src)
    if err != nil {
        return Errors.NewWarning(err.Error())
    }
    if !sfi.Mode().IsRegular() {
        // cannot copy non-regular files (e.g., directories,
        // symlinks, devices, etc.)
        return Errors.NewWarning("CopyFile: non-regular source file ", sfi.Name(), "(", sfi.Mode().String(), ")")
    }
    dfi, err := os.Stat(dst)
    if err != nil {
        if !os.IsNotExist(err) {
            return Errors.NewWarning(err.Error())
        }
    } else {
        if !(dfi.Mode().IsRegular()) {
            return Errors.NewWarning("CopyFile: non-regular destination file ", dfi.Name(), "(", dfi.Mode().String(), ")")
        }
        if os.SameFile(sfi, dfi) {
            return Errors.None()
        }
    }
    if err = os.Link(src, dst); err == nil {
        return Errors.None()
    }

    return copyFileContents(src, dst)
}


// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (Errors.Error) {
    in, err := os.Open(src)
    if err != nil {
        return Errors.NewFatal("boo", err.Error())
    }
    defer in.Close()

    // Create destination folder
    splitPath := Helpers.Split(dst, "\\")
    folder := Helpers.Join(splitPath[:len(splitPath) - 1], "\\")
    err1 := os.MkdirAll(folder, 0777)
    if err1 != nil {
        return Errors.NewFatal(err1.Error())
    }

    out, err := os.Create(dst)
    if err != nil {
        return Errors.NewFatal("OKAY",err.Error())
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, in); err != nil {
        return Errors.None()
    }
    err = out.Sync()
    return Errors.None()
}


/**
 * Name.........: FileExists
 * Parameters...: path (string) - path to the file to check
 * Return.......: bool - true if the file exists
 * Description..: Checks if a file exists
 */
func FileExists (name string) bool {
  _, result := os.Stat(name)

  if result != nil {
    if os.IsNotExist(result) {
      return false
    }
  }
  return true
}


/**
  * Name.........: EmptyDir
  * Parameters...: dir (string) - The directory to empty
  * Description..: Removes all files in a directory
  */
func EmptyDir(dir string) {
    list, err := ioutil.ReadDir(dir)
    if err != nil {
        return
    }

    // Loop through all files/folders in the directory
    for _, item := range list {
        if item.IsDir() {
            // Delete all items in that folder
            EmptyDir(dir + "\\" + item.Name())

        } else {
            // Delete the file
            os.Remove(dir + "\\" + item.Name())
        }
    }

    // Remove the root directory
    os.Remove(dir)
}


/**
  * Name.........: CollapseDirectory
  * Parameters...: dir (string) - directory to collapse
  *                recursive (bool) - if true then it will get files from subfolders
  * Return.......: []os.FileInfo - all the files found
  * Description..: Finds files in a directory, and possible subdirectories
  */
func CollapseDirectory(root string, dir string, recursive bool) ([]DataTypes.FileToCompile) {
    files := []DataTypes.FileToCompile{}

    if len(dir) > 0 && dir[:1] == "\\" {
        dir = dir[1:]
    }

    list, err := ioutil.ReadDir(root + "\\" + dir)
    if err != nil {
        return files
    }

    if len(dir) > 1 && dir[:1] == "." {
        return files
    }

    for _, item := range list {
        if item.IsDir() {
            if recursive {
                files = append(files, CollapseDirectory(root, dir + "\\" + item.Name(), true)...)
            }
        } else {
            //fmt.Println("Adding: ", item.Name())
            files = append(files, DataTypes.FileToCompile{Directory:dir,Info:item})
        }
    }

    return append([]DataTypes.FileToCompile{}, files...)
}
