package http

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-faster/errors"
	"github.com/gorilla/mux"
	"github.com/gotd/contrib/http_io"
	"github.com/gotd/contrib/partio"
	"github.com/gotd/contrib/tg_io"
	"github.com/gotd/td/tg"

	"github.com/iyear/tdl/core/logctx"
	"github.com/iyear/tdl/core/tmedia"
	"github.com/shivamhw/content-pirate/pkg/log"
	"github.com/shivamhw/content-pirate/pkg/telegram"
)

type media struct {
	*tmedia.Media
	MIME string
}

var (
	tClient *telegram.Telegram
)

func LoginTelegram(phone string, otp string) error {
	ctx := context.Background()
	t, err := telegram.NewTelegram(ctx, &telegram.UserData{
		PhoneNumber: phone,
	})
	if err != nil {
		return fmt.Errorf("failed to create telegram client: %w", err)
	}

	err = t.Login(&telegram.LoginOpts{
		Phone: phone,
		Otp:   otp,
	}, false)
	if err != nil {
		return fmt.Errorf("failed to login to telegram: %w", err)
	}
	tClient = t
	log.Infof("Telegram client initialized successfully with phone: %s", phone)
	return nil
}

func InitTelegramClient(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number is required to initialize telegram client")
	}

	ctx := context.Background()
	t, err := telegram.NewTelegram(ctx, &telegram.UserData{
		PhoneNumber: phone,
	})
	if err != nil {
		return fmt.Errorf("failed to create telegram client: %w", err)
	}

	tClient = t
	log.Infof("Telegram client initialized successfully with phone: %s", phone)
	return nil
}

func StreamHandler(ctx context.Context) http.Handler {

	cache := &sync.Map{} // map[string]*media
	//router.Handle("/{peer}/{message:[0-9]+}"
	return handler(func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		peer := vars["peer"]
		messageStr := vars["message"]

		var item *media
		if t, ok := cache.Load(peer + messageStr); ok {
			item = t.(*media)
		} else {
			message, err := strconv.Atoi(messageStr)
			if err != nil {
				return errors.Wrap(err, "invalid message id")
			}
			msg, err := tClient.GetSingleMessage(message, peer)
			item, err = convItem(msg)
			if err != nil {
				return errors.Wrap(err, "convItem")
			}

			cache.Store(peer+messageStr, item)
		}

		api := tClient.GetClient()
		

		w.Header().Set("Connection", "keep-alive")

		u := partio.NewStreamer(
			tg_io.NewDownloader(api.API()).ChunkSource(item.Size, item.InputFileLoc),
			int64(512*1024))

		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, item.Name))
		http_io.NewHandler(u, item.Size).
			WithContentType(item.MIME).
			WithLog(logctx.From(ctx).Named("serve")).
			ServeHTTP(w, r)
		return nil
	})
}

func handler(h func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}

func convItem(msg *tg.Message) (*media, error) {
	md, ok := tmedia.GetMedia(msg)
	if !ok {
		return nil, errors.New("message is not a media")
	}

	mime := ""
	switch m := msg.Media.(type) {
	case *tg.MessageMediaDocument:
		doc, ok := m.Document.AsNotEmpty()
		if !ok {
			return nil, errors.New("document is empty")
		}
		mime = doc.MimeType
	case *tg.MessageMediaPhoto:
		mime = "image/jpeg"
	}

	return &media{
		Media: md,
		MIME:  mime,
	}, nil
}
