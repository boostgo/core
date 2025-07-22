package build

import (
	"github.com/boostgo/core/appx"
	"github.com/boostgo/core/authx"
	"github.com/boostgo/core/configx"
	"github.com/boostgo/core/convert"
	"github.com/boostgo/core/echox"
	"github.com/boostgo/core/errorx"
	"github.com/boostgo/core/fsx"
	"github.com/boostgo/core/grpcx"
	"github.com/boostgo/core/httpx"
	"github.com/boostgo/core/kafkax"
	"github.com/boostgo/core/log"
	"github.com/boostgo/core/log/logx"
	"github.com/boostgo/core/mathx"
	"github.com/boostgo/core/mongox"
	"github.com/boostgo/core/orderedmap"
	"github.com/boostgo/core/pagex"
	"github.com/boostgo/core/queuex"
	"github.com/boostgo/core/redis"
	"github.com/boostgo/core/reflectx"
	"github.com/boostgo/core/requests"
	"github.com/boostgo/core/retry"
	"github.com/boostgo/core/samples"
	"github.com/boostgo/core/semaphore"
	"github.com/boostgo/core/sql"
	"github.com/boostgo/core/storage"
	"github.com/boostgo/core/trace"
	"github.com/boostgo/core/translate"
	"github.com/boostgo/core/translit"
	"github.com/boostgo/core/tsv"
	"github.com/boostgo/core/validator"
	"github.com/boostgo/core/worker"
)

func Build() {
	logx.Pretty()
	log.Info().Msg("log +")

	trace.IAmMaster(true)
	log.Info().Msg("trace +")

	convert.String(1)
	log.Info().Msg("convert +")

	_ = errorx.New("test error")
	log.Info().Msg("errorx +")

	_ = appx.Context()
	log.Info().Msg("appx +")

	_ = validator.New()
	log.Info().Msg("validator +")

	_ = tsv.ErrNoRecords
	log.Info().Msg("tsv +")

	_ = sql.ErrConnectionIsNotShard
	log.Info().Msg("sql +")

	_ = redis.ErrClientAddressEmpty
	log.Info().Msg("redis +")

	_ = storage.ErrConnNotSelected
	log.Info().Msg("storage +")

	_ = retry.ErrMaxRetriesExceeded
	log.Info().Msg("retry +")

	_ = worker.ErrLocked
	log.Info().Msg("worker +")

	_ = grpcx.Registry(nil)
	log.Info().Msg("grpc +")

	_ = fsx.AnyFileExist()
	log.Info().Msg("fsx +")

	_ = httpx.ContentTypeJSON
	log.Info().Msg("httpx +")

	_ = mathx.Abs(0)
	log.Info().Msg("mathx +")

	_ = orderedmap.NewOrderedMap[string, string]()
	log.Info().Msg("orderedmap +")

	_ = pagex.Pagination{}
	log.Info().Msg("pagex +")

	_ = queuex.Config{}
	log.Info().Msg("queuex +")

	_ = reflectx.IsPointer(nil)
	log.Info().Msg("reflectx +")

	_ = authx.ErrNoToken
	log.Info().Msg("authx +")

	_ = configx.GetString("")
	log.Info().Msg("configx +")

	_ = echox.RouterGroup{}
	log.Info().Msg("echox +")

	_ = kafkax.ErrConnect
	log.Info().Msg("kafkax +")

	_ = requests.ErrBytesWriterWrite
	log.Info().Msg("requests +")

	_ = translate.ErrKeyNotFound
	log.Info().Msg("translate +")

	_ = samples.Server{}
	log.Info().Msg("samples +")

	_ = semaphore.Semaphore{}
	log.Info().Msg("semaphore +")

	_ = translit.LatinByCyrillic("")
	log.Info().Msg("translit +")

	_ = mongox.ErrConcernReadUnsupported
	log.Info().Msg("mongox +")
}
