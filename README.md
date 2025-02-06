# Gofali

A simple URL shortener written in Go.

![Gofali img](https://github.com/user-attachments/assets/8d7281ef-7cb9-4a4e-8089-e63c7ab5ea33)

## Run

1. Grab the latest release
2. Create `config.ini`
3. Run `gofali-<platform> -config config.ini`

## Configuration

You can configure the application using flags or a `config.ini` file. Below are the available options:

### Flags

| Flag              | Default Value                              | Description                               |
|-------------------|--------------------------------------------|-------------------------------------------|
| `-host`           | `127.0.0.1`                               | Server host                                |
| `-port`           | `8080`                                    | Server port                                |
| `-mysql-connect`  | `gofali:gofali@tcp(127.0.0.1:3306)/gofali` | MySQL connection string                   |
| `-admin-path`     | `admin`                                   | Admin path (`example.com/@admin/`)         |
| `-root-redirect`  | (empty)                                   | Redirect path for the root URL            |
| `-slug-length`    | `8`                                       | Length of the generated slug              |

### Example `config.ini`

```ini
host = 127.0.0.1
port = 8080

mysql-connect = gofali:gofali@tcp(127.0.0.1:3306)/gofali

admin-path = BnTiC5DZ3L
root-redirect = https://google.com
slug-length = 16
```
