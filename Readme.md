## Create go project that accessing the database using postgreSql

## How to run the project

### Prerequisites

- set user and password by env using terminal or cmd

- for linux or mac

  ```bash
  $ export DBUSER=postgres
  $ export DBPASS=postgres
  ```

  

- for windows

  ```bash
  C:\Users\you\data-access> set DBUSER=postgres
  C:\Users\you\data-access> set DBPASS=postgres
  ```

### Step by Step to Running

- for linux or mac => `$ ./data-access`
- for windows => `data-access.exe`
- or open terminal and run `go run .`

### Features:

- connect to psql database

  ```go
  // open connection
  var err error
  db, err = sql.Open("postgres", psqlInfo)
  checkError(err)
  defer db.Close()
  ```

  

- get realtime error from database

  ```go
  pingErr := db.Ping()
  checkError(pingErr)
  ```

  

- get albums by artist name

  ```go
  func albumsByArtist(name string) ([]Album, error)
  ```

  

- get album by id

  ```go
  func getAlbumById(id int64) (Album, error)
  ```

  

- add new album

  ```go
  func addAlbum(album Album) (string, error)
  ```

  

- get all album

  ```go
  func getAlbums() ([]Album, error)
  ```

  

- delete album by id

  ```go
  func deleteAlbumById(id int64) (string, error)
  ```

  

- update album by id

  ```go
  func updateAlbumById(id int64, album Album) (Album, error)
  ```

  
