package entity

type stages struct {
	extractors        []string
	isQuality         bool
	isCookie          bool
	isMarkWatched     bool
	isEmbededCharters bool
}

func (f stages) GetIsQuality() bool {
	return f.isQuality
}

func (f stages) GetExtractors() []string {
	return f.extractors
}

func (f stages) GetIsCookie() bool {
	return f.isCookie
}

func (f stages) GetIsMarkWatched() bool {
	return f.isMarkWatched
}

func (f stages) GetIsEmbededCharters() bool {
	return f.isEmbededCharters
}

type task struct {
	maxQuality int32
	link       string
	Stage      []stages
}

func (t task) GetQuality() int32 {
	return t.maxQuality
}

func (t task) GetLink() string {
	return t.link
}

var (
	extrctrs = []string{"", "youtube:player-client=tv_embedded", "youtube:player_client=default,ios", "youtube:player-client=default,-tv,web_safari,web_embedded"}

	stgs = []stages{
		{
			isQuality:         true,
			isCookie:          true,
			isEmbededCharters: true,
			isMarkWatched:     true,
			extractors:        extrctrs,
		},
		{
			isQuality:         true,
			isCookie:          true,
			isEmbededCharters: true,
			isMarkWatched:     false,
			extractors:        extrctrs,
		},
		{
			isQuality:         true,
			isCookie:          true,
			isEmbededCharters: false,
			isMarkWatched:     true,
			extractors:        extrctrs,
		},
		{
			isQuality:         true,
			isCookie:          true,
			isEmbededCharters: false,
			isMarkWatched:     false,
			extractors:        extrctrs,
		},
		{
			isQuality:         true,
			isCookie:          false,
			isEmbededCharters: true,
			isMarkWatched:     true,
			extractors:        extrctrs,
		},
		{
			isQuality:         true,
			isCookie:          false,
			isEmbededCharters: true,
			isMarkWatched:     false,
			extractors:        extrctrs,
		},
		{
			isQuality:         true,
			isCookie:          false,
			isEmbededCharters: false,
			isMarkWatched:     true,
			extractors:        extrctrs,
		},
		{
			isQuality:         true,
			isCookie:          false,
			isEmbededCharters: false,
			isMarkWatched:     false,
			extractors:        extrctrs,
		},

		{
			isQuality:         false,
			isCookie:          true,
			isEmbededCharters: true,
			isMarkWatched:     true,
			extractors:        extrctrs,
		},
		{
			isQuality:         false,
			isCookie:          true,
			isEmbededCharters: true,
			isMarkWatched:     false,
			extractors:        extrctrs,
		},
		{
			isQuality:         false,
			isCookie:          true,
			isEmbededCharters: false,
			isMarkWatched:     true,
			extractors:        extrctrs,
		},
		{
			isQuality:         false,
			isCookie:          true,
			isEmbededCharters: false,
			isMarkWatched:     false,
			extractors:        extrctrs,
		},
		{
			isQuality:         false,
			isCookie:          false,
			isEmbededCharters: true,
			isMarkWatched:     true,
			extractors:        extrctrs,
		},
		{
			isQuality:         false,
			isCookie:          false,
			isEmbededCharters: true,
			isMarkWatched:     false,
			extractors:        extrctrs,
		},
		{
			isQuality:         false,
			isCookie:          false,
			isEmbededCharters: false,
			isMarkWatched:     true,
			extractors:        extrctrs,
		},
		{
			isQuality:         false,
			isCookie:          false,
			isEmbededCharters: false,
			isMarkWatched:     false,
			extractors:        extrctrs,
		},
		{
			isQuality:         true,
			isCookie:          true,
			isEmbededCharters: true,
			isMarkWatched:     true,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         true,
			isCookie:          true,
			isEmbededCharters: true,
			isMarkWatched:     false,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         true,
			isCookie:          true,
			isEmbededCharters: false,
			isMarkWatched:     true,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         true,
			isCookie:          true,
			isEmbededCharters: false,
			isMarkWatched:     false,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         true,
			isCookie:          false,
			isEmbededCharters: true,
			isMarkWatched:     true,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         true,
			isCookie:          false,
			isEmbededCharters: true,
			isMarkWatched:     false,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         true,
			isCookie:          false,
			isEmbededCharters: false,
			isMarkWatched:     true,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         true,
			isCookie:          false,
			isEmbededCharters: false,
			isMarkWatched:     false,
			extractors:        []string{"youtube:formats=missing_pot"},
		},

		{
			isQuality:         false,
			isCookie:          true,
			isEmbededCharters: true,
			isMarkWatched:     true,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         false,
			isCookie:          true,
			isEmbededCharters: true,
			isMarkWatched:     false,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         false,
			isCookie:          true,
			isEmbededCharters: false,
			isMarkWatched:     true,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         false,
			isCookie:          true,
			isEmbededCharters: false,
			isMarkWatched:     false,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         false,
			isCookie:          false,
			isEmbededCharters: true,
			isMarkWatched:     true,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         false,
			isCookie:          false,
			isEmbededCharters: true,
			isMarkWatched:     false,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         false,
			isCookie:          false,
			isEmbededCharters: false,
			isMarkWatched:     true,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
		{
			isQuality:         false,
			isCookie:          false,
			isEmbededCharters: false,
			isMarkWatched:     false,
			extractors:        []string{"youtube:formats=missing_pot"},
		},
	}
)

func NewTask(quality int32, link string) *task {
	return &task{
		maxQuality: quality,
		link:       link,
		Stage:      stgs,
	}
}
