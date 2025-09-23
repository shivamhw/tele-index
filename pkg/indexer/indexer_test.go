package indexer

import (
	"fmt"
	"testing"

	"github.com/shivamhw/tele-index/internal/models"
)

var (
	items = []models.ItemAddRequest{
		{Id: 102550, FileName: "The Isle(2000)@msm.mkv", ChatId: 2714472314, From: 7484725854, Size: 1411758877},
		{Id: 102549, FileName: "[ME] Alice In Borderland (2022) S02E01 Dua Aud 690MB.mkv", ChatId: 2714472314, From: 7484725854, Size: 729448367},
		{Id: 102548, FileName: "(@UCMOVIE) The.Isle.2019.720p.WEB-DL.MkvCage.mkv", ChatId: 2714472314, From: 7484725854, Size: 841013894},
		{Id: 102547, FileName: "@MM_New Commando 2 [2017] Tamil Dubbed Dvdrip x264 700MB.mkv", ChatId: 2714472314, From: 7484725854, Size: 776028913},
		{Id: 102546, FileName: "The_Kashmir_Files_2022_1080p_10bit_ZEE5_WEBRip_x265_HEVC_Hindi_DDP.mkv", ChatId: 2714472314, From: 7484725854, Size: 1932134806},
		{Id: 102545, FileName: "@ADrama_Lovers- Alice.In.Borderland.S02E01.NF.x264.540p.mkv", ChatId: 2714472314, From: 7484725854, Size: 283841578},
		{Id: 102544, FileName: "The.Crown.S01E03.1080p.10bit.WEBRip.6CH.x265.HEVC-PSA.mkv", ChatId: 2714472314, From: 7484725854, Size: 750164104},
		{Id: 102543, FileName: "@RickyChannel Commando.2.2017.720p.Dvdrip.x264-NBY.mkv", ChatId: 2714472314, From: 7484725854, Size: 994512914},
		{Id: 102542, FileName: "The_Kashmir_Files_2022_1080p_10bit_Z5_WR_DDP5_1_NVENC_Bunny_mkv.mkv", ChatId: 2714472314, From: 7484725854, Size: 1692496992},
		{Id: 102541, FileName: "@TM_LMO Nenjil Thunivirundhal (2017) Tamil HDRip.mkv", ChatId: 2714472314, From: 7484725854, Size: 735516498},
		{Id: 102540, FileName: "Delhi.Crime.S02E04.720p.10bit.NF.WEBRip.HIN.AAC5.1.x265.HEVC.mkv", ChatId: 2714472314, From: 7484725854, Size: 186897989},
		{Id: 102539, FileName: "Fear_the_Walking_Dead_S01_E02_WebRip_Hindi_5_1_+_English_5_1_720p.mkv", ChatId: 2714472314, From: 7484725854, Size: 461565990},
		{Id: 102538, FileName: "Immini_Nalloraal_2005_1080p_8bit_MMAX_WEB_DL_AAC_2_0_x265_AVK.mkv", ChatId: 2714472314, From: 7484725854, Size: 1541307014},
		{Id: 102537, FileName: "[CV] Thalaivettiyaan Paalayam S01E06 Inimey pechukku ida.mkv", ChatId: 2714472314, From: 7484725854, Size: 277108126},
		{Id: 102536, FileName: "Surviving Summer S02E05 Hin-Eng (720p WEBRip 10bit x265) [PM.mkv", ChatId: 2714472314, From: 7484725854, Size: 220035794},
		{Id: 102535, FileName: "[CV] Thalaivettiyaan Paalayam S01E08 Sound Amma 720p.mkv", ChatId: 2714472314, From: 7484725854, Size: 317378274},
		{Id: 102534, FileName: "Vikram 2022 UNCUT 1080p 10bit HEVC HDRip ORG. [Hindi DD 5.1 .mkv", ChatId: 2714472314, From: 7484725854, Size: 2691229264},
		{Id: 102533, FileName: "Surviving Summer S02E05 Hin-Eng (1080p WEBRip 10bit x265 C0S.mkv", ChatId: 2714472314, From: 7484725854, Size: 581834065},
		{Id: 102532, FileName: "Sherlock_Holmes_A_Game_Of_Shadows_2011_1080p_BluRay_x265_RARBG.mp4", ChatId: 2714472314, From: 7484725854, Size: 2154220312},
		{Id: 102531, FileName: "Alice_in_Borderland_S02E01_1080p_NF_WEB_DL_MULTI_DDP5_1_Atmos_HDR.mkv", ChatId: 2714472314, From: 7484725854, Size: 1182298852},
		{Id: 102530, FileName: "Alice_in_Borderland_S02E01_1080p_NF_WEB_DL_DDP5_1_Atmos_HEVC_Saon.mkv", ChatId: 2714472314, From: 7484725854, Size: 1467266706},
	}
)

func setup() *Indexer {
	indx, err := NewIndexer(&IndexerOpts{
		IndexSchema: "./index.json",
		IndexPath:   "./testIdx",
	})
	if err != nil {
		panic(err)
	}
	return indx
}


func loadData(indx *Indexer) error {
	for _, i := range items {
		it := i.GetItem()
		if err := indx.Index(it); err != nil {
		return err
		}
	}
	return nil
}

func TestIndex(t *testing.T) {
	indx := setup()
	defer indx.Close()
	if err := loadData(indx); err != nil {
		t.Fatal(err)
	}
	d, err := indx.GetTotalDocs()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != int(d) {
		t.Fatalf("mismatch: items: %d, idx: %d", len(items), d)
	}
	fmt.Printf("matched: items: %d, idx: %d", len(items), d)
}

func TestSearch(t *testing.T) {
	indx := setup()
	defer indx.Close()
	if err := loadData(indx); err != nil {
		t.Fatal(err)
	}
	res, err := indx.Search("s02e01", 10, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 4 {
		t.Fatalf("mismatch, have %d want %d", len(res), 4)
	}
	fmt.Printf("match, have %d want %d ", len(res), 4)
}

func TestBulkImport(t *testing.T){
	indx := setup()
	defer indx.Close()
	var its []*models.Item
	for _, it := range items {
		t := it.GetItem()
		its = append(its, &t)
	}

	if err := indx.BulkImport(its); err != nil {
		t.Fatal(err)
	}
	d, err := indx.GetTotalDocs()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != int(d) {
		t.Fatalf("mismatch, want %d have %d", len(items), d)
	}
}