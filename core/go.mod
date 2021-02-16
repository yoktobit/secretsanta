module github.com/yoktobit/secretsanta

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/atkinsonbg/go-gmux-db-testcontainers v0.0.0-20210123160844-a46ed7ad99ed
	github.com/docker/go-connections v0.4.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/google/uuid v1.1.2 // indirect
	github.com/google/wire v0.5.0
	github.com/jackc/pgx/v4 v4.9.0
	github.com/lithammer/shortuuid v3.0.0+incompatible
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/sirupsen/logrus v1.4.2
	github.com/testcontainers/testcontainers-go v0.9.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gorm.io/driver/postgres v1.0.5
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.8
)

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20190813064441-fde4db37ae7a
