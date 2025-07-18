module github.com/RenterRus/dwld-downloader

go 1.24.4

require (
	github.com/AlekSi/pointer v1.2.0
	github.com/bradfitz/gomemcache v0.0.0-20250403215159-8d39553ac7cf
	github.com/go-playground/validator/v10 v10.27.0
	github.com/lrstanley/go-ytdlp v1.1.0
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/samber/lo v1.51.0
	github.com/spf13/viper v1.20.1
	google.golang.org/grpc v1.74.0
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/RenterRus/dwld-ftp-sender v0.0.0-20250717203705-b097892624c5

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/crypto v0.40.0 // indirect
)

require (
	github.com/ProtonMail/go-crypto v1.3.0 // indirect
	github.com/RenterRus/dwld-bot v0.0.0-20250717004310-b68972ace5d3
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/redis/go-redis/v9 v9.11.0
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
)

replace github.com/RenterRus/dwld-ftp-sender => github.com/RenterRus/dwld-ftp-sender v0.0.0-20250717205716-0c8f04bc8617
