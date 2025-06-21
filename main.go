package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/Noooste/azuretls-client"
	"github.com/gurkankaymak/gosoup"
	"github.com/tealeg/xlsx/v3"
	"strings"
	"time"
)

type flagDTO struct {
	saveAll           bool
	tenderPages       int
	orderPages        int
	appendAll         bool
	tenderOldFileName string
	ordersOldFileName string
}

// Dostawa wyposażenia diagnostycznego urządzeń elektronicznych – pomocy dydaktycznych niezbędnych do realizacji kursów w ramach projektu „Rozwój CKZ w Nowym Sączu”. Postępowanie 3.
// Dostawa wyposażenia diagnostycznego urządzeń elektronicznych – pomocy dydaktycznych niezbędnych do realizacji kursów w ramach projektu „Rozwój CKZ w Nowym Sączu”. Postępowanie 3.
// https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/view/3108196/dostawa-wyposazenia-diagnostycznego-urzadzen-elektronicznych-pomocy-dydaktycznych-niezbednych-do-realizacji-kursow-w-ramach?_7_WAR_organizationnoticeportlet_redirect=https%3A%2F%2Foneplace.marketplanet.pl%2Fzapytania-ofertowe-przetargi%2F-%2Frfp%2Fcat%3F_7_WAR_organizationnoticeportlet_
// cur%3D1
// &_7_WAR_organizationnoticeportlet_friendlyUrl=dostawa-wyposazenia-diagnostycznego-urzadzen-elektronicznych-pomocy-dydaktycznych-niezbednych-do-realizacji-kursow-w-ramach&_7_WAR_organizationnoticeportlet
// _noticeId=3108196
// https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/view/3108196/dostawa-wyposazenia-diagnostycznego-urzadzen-elektronicznych-pomocy-dydaktycznych-niezbednych-do-realizacji-kursow-w-ramach?_7_WAR_organizationnoticeportlet_redirect=https%3A%2F%2Foneplace.marketplanet.pl%2Fzapytania-ofertowe-przetargi%2F-%2Frfp%2Fcat%3F_7_WAR_organizationnoticeportlet_
// cur%3D31
// &_7_WAR_organizationnoticeportlet_friendlyUrl=dostawa-wyposazenia-diagnostycznego-urzadzen-elektronicznych-pomocy-dydaktycznych-niezbednych-do-realizacji-kursow-w-ramach&_7_WAR_organizationnoticeportlet
// _noticeId=3108196
func newFlagDTO() *flagDTO {
	saveAll := flag.Bool("saveAll", true, "save all tenders to excel")
	//max 1048567/12=87381, found 1090*12 =>1000
	tenderPages := flag.Int("tenderPages", 1000, "number of tender Pages to get")
	//max 1048567/50=20971, found 18917:50=378 => 300 => server error response of max=1000
	//max 1000/50=200!!!
	orderPages := flag.Int("orderPages", 200, "number of order Pages to get")
	appendAll := flag.Bool("appendAll", false, "append old tenders to new all")
	tenderOldFileName := flag.String("tenderOldFileName", "przetargi_all.xlsx", "tender old file name")
	ordersOldFileName := flag.String("orderOldFileName", "oferty_all.xlsx", "order old file name")
	flags := flagDTO{
		saveAll:           *saveAll,
		tenderPages:       *tenderPages,
		orderPages:        *orderPages,
		appendAll:         *appendAll,
		tenderOldFileName: *tenderOldFileName,
		ordersOldFileName: *ordersOldFileName}
	return &flags
}

type tenderDTO struct {
	name string
	href string
	date string
	id   string
	isIT bool
}

func newTenderDTO(name string, href string, date string, id string) *tenderDTO {
	lowerName := strings.ToLower(name)
	isIT := strings.Contains(lowerName, "oprogramowani") ||
		strings.Contains(lowerName, " it ") ||
		strings.Contains(lowerName, "rozwój i utrzymanie systemu") ||
		strings.Contains(lowerName, "aplikacj")
	p := tenderDTO{name: name, href: href, date: date, id: id, isIT: isIT}
	return &p
}

// [{"objectId":"ocds-148610-c02258fe-182b-4522-8016-10439c523c4e", "title":"„Świadczenie usług transportowych w zakresie przewozów uczniów niepełnosprawnych \ndo szkół i placówek oraz sprawowanie opieki nad uczniami podczas przewozów w roku szkolnym 2025/2026”.","organizationId":"8037","organizationName":"GMINNA ADMINISTRACJA OŚWIATY","organizationPartName":null,"organizationCity":"Tuchów","organizationProvince":"małopolskie","bzpNumber":"2025/BZP 00284511/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:59:06.709Z"},{"objectId":"ocds-148610-74cd5c84-2e9e-4d24-bf40-de65d6c9bf8c","title":"Modernizacja ogólnodostępnego placu zabaw w Koszycach","organizationId":"2713","organizationName":"GMINA KOSZYCE","organizationPartName":null,"organizationCity":"Koszyce","organizationProvince":"małopolskie","bzpNumber":"2025/BZP 00284506/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-04T10:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:58:12.729Z"},{"objectId":"ocds-148610-00528711-874e-40d6-922c-22c709cfb635","title":"Modernizacja (przebudowa) drogi dojazdowej do gruntów rolnych \nw miejscowości Kocudza Trzecia na działkach o nr ewid. 1039, 1058 i 1095.","organizationId":"14399","organizationName":"Gmina Dzwola","organizationPartName":null,"organizationCity":"Dzwola","organizationProvince":"lubelskie","bzpNumber":"2025/BZP 00284502/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:57:36.843Z"},{"objectId":"ocds-148610-ba280dd4-0ad1-43f8-bd81-1723558c7390","title":"Wyposażenie pracowni językowej w RCEZ w Nisku w ramach projektu pn.: \"Wzmocnienie potencjału szkół zawodowych w Powiecie Niżańskim\"","organizationId":"3463","organizationName":"Powiat Niżański","organizationPartName":null,"organizationCity":"Nisko","organizationProvince":"podkarpackie","bzpNumber":"2025/BZP 00284501/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T09:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:57:31.5Z"},{"objectId":"ocds-148610-a1ec52c4-1be2-4303-8c63-1b993fd0b99e","title":"Budowa odnawialnych źródeł energii w gminie Radziechowy-Wieprz","organizationId":"7417","organizationName":"Gmina Radziechowy-Wieprz","organizationPartName":null,"organizationCity":"Radziechowy","organizationProvince":"śląskie","bzpNumber":"2025/BZP 00284498/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-09T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:57:05.174Z"},{"objectId":"ocds-148610-ab627c45-1f3e-4961-bbba-a2594367cc35","title":"Zakup i dostawa aparatu USG oraz sprzętu i wyposażenia medycznego","organizationId":"82593","organizationName":"GMINNY OŚRODEK ZDROWIA W WODZISŁAWIU","organizationPartName":null,"organizationCity":"Wodzisław","organizationProvince":"świętokrzyskie","bzpNumber":"2025/BZP 00284494/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-26T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:56:32.988Z"},{"objectId":"ocds-148610-812a35a6-46f0-42c7-b6a8-816e43e921cf","title":"Dostawa urządzeń sieciowych","organizationId":"9169","organizationName":"TEATR NARODOWY","organizationPartName":null,"organizationCity":"Warszawa","organizationProvince":"mazowieckie","bzpNumber":"2025/BZP 00284487/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:55:25.253Z"},{"objectId":"ocds-148610-9d84a1ff-097f-479b-81f4-dd6d0732b023","title":"Remont i modernizacja kompleksu sportowego „Moje Boisko - Orlik 2012” w Wielgolesie, gm. Latowicz","organizationId":"6797","organizationName":"Gmina Latowicz","organizationPartName":null,"organizationCity":"Latowicz","organizationProvince":"mazowieckie","bzpNumber":"2025/BZP 00284484/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:54:27.91Z"},{"objectId":"ocds-148610-24cd4e84-7cbe-4ff6-b14e-8e2f6e945a6d","title":"Budowa utwardzonego pobocza w ciągu drogi powiatowej nr 2162K w miejscowości Wężerów","organizationId":"9513","organizationName":"ZARZĄD DRÓG POWIATU KRAKOWSKIEGO","organizationPartName":null,"organizationCity":"Batowice","organizationProvince":"małopolskie","bzpNumber":"2025/BZP 00284465/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-04T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:51:34.216Z"},{"objectId":"ocds-148610-f2e81220-9c1e-44f1-8963-6a2dbd9b32b3","title":"Dowóz uczniów niepełnosprawnych z terenu Gminy Goleszów do szkół, ośrodków i placówek oświatowych.","organizationId":"1429","organizationName":"GMINA GOLESZÓW","organizationPartName":null,"organizationCity":"Goleszów","organizationProvince":"śląskie","bzpNumber":"2025/BZP 00284431/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-26T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:44:03.546Z"},{"objectId":"ocds-148610-347655b3-bc3d-4b86-8b7f-721f3663f684","title":"Dostawa, skonfigurowanie i uruchomienie zestawu magnetometrów do pomiarów parametrów magnetycznych elementów ferromagnetycznych wraz z akcesoriami","organizationId":"1656","organizationName":"Uniwersytet Radomski im. Kazimierza Pułaskiego","organizationPartName":null,"organizationCity":"Radom","organizationProvince":"mazowieckie","bzpNumber":"2025/BZP 00284427/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T09:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:43:32.541Z"},{"objectId":"ocds-148610-d08bb451-63f2-4dc6-93c6-652df3a21c0d","title":"Odbudowa drogi gminnej nr 510562K – ul. Kacza w Jawiszowicach w km 0+630 – 0+869 wraz z odbudową mostu w km 0+865","organizationId":"1136","organizationName":"GMINA BRZESZCZE","organizationPartName":null,"organizationCity":"Brzeszcze","organizationProvince":"małopolskie","bzpNumber":"2025/BZP 00284424/01","tenderType":"1.1.2","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-04T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:43:04.76Z"},{"objectId":"ocds-148610-7a79759d-2a26-429a-aa61-341e3881a158","title":"Remont toalet w Szkole Podstawowej w Jeziorzanach.","organizationId":"11139","organizationName":"Gmina Liszki","organizationPartName":null,"organizationCity":"Liszki","organizationProvince":"małopolskie","bzpNumber":"2025/BZP 00284417/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-02T10:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:40:20.102Z"},{"objectId":"ocds-148610-0983438b-7ac3-4a0c-b5f0-de1cb01e647f","title":"Dostawa produktów farmaceutycznych z podziałem na 14 części dla Powiatowego Szpitala im. Władysława Biegańskiego w Iławie","organizationId":"1483","organizationName":"Powiatowy Szpital im. Władysława Biegańskiego w Iławie","organizationPartName":null,"organizationCity":"Iława","organizationProvince":"warmińsko-mazurskie","bzpNumber":"2025/BZP 00284403/01","tenderType":"1.1.2","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-26T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:36:33.775Z"},{"objectId":"ocds-148610-250b3815-5f6c-42c4-aaf7-d792fd4b58bc","title":"„Zorganizowanie i przeprowadzenie części praktycznej kursu prawa jazdy kat. B dla uczniów ZSCKR w Kamieniu Małym\"","organizationId":"12616","organizationName":"ZESPÓŁ SZKÓŁ CENTRUM KSZTAŁCENIA ROLNICZEGO W KAMIENIU MAŁYM","organizationPartName":null,"organizationCity":"Kamień Mały","organizationProvince":"lubuskie","bzpNumber":"2025/BZP 00284402/01","tenderType":"1.1.2","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:36:32.489Z"},{"objectId":"ocds-148610-d87c63bd-ab48-46aa-91c4-ac7d7516676f","title":"„Remont drogi powiatowej nr 2521D Czernica”","organizationId":"8501","organizationName":"Powiat Karkonoski - Starostwo Powiatowe w Jeleniej Górze","organizationPartName":null,"organizationCity":"Jelenia Góra","organizationProvince":"dolnośląskie","bzpNumber":"2025/BZP 00284394/01","tenderType":"1.1.2","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:34:41.537Z"},{"objectId":"ocds-148610-23000103-0a30-454d-a80c-dfd16e095ef4","title":"Zimowe utrzymanie dróg gminnych i wewnętrznych w sezonie 2025/2026 na terenie gminy Lipowa.","organizationId":"4931","organizationName":"Gmina Lipowa","organizationPartName":null,"organizationCity":"Lipowa","organizationProvince":"śląskie","bzpNumber":"2025/BZP 00284387/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-30T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:32:39.253Z"},{"objectId":"ocds-148610-ac52fef0-4673-43df-a7fb-e109236db2da","title":"Przygotowanie, dostarczanie i wydawanie posiłków uczniom Szkoły Podstawowej nr 2 im. Henryka Sienkiewicza w Wieluniu","organizationId":"25012","organizationName":"SZKOŁA PODSTAWOWA NR 2 IM. HENRYKA SIENKIEWICZA W WIELUNIU","organizationPartName":null,"organizationCity":"Wieluń","organizationProvince":"łódzkie","bzpNumber":"2025/BZP 00284375/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-30T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:30:37.848Z"},{"objectId":"ocds-148610-1d8a2374-b5e0-4fae-a2a9-4a96a08ea776","title":"Ustawienie toalet przenośnych na terenie miasta Krakowa wraz z ich serwisowaniem w latach 2025-2028","organizationId":"3537","organizationName":"Zarząd Infrastruktury Wodnej","organizationPartName":null,"organizationCity":"Kraków","organizationProvince":"małopolskie","bzpNumber":"2025/BZP 00284368/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-02T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:29:39.352Z"},{"objectId":"ocds-148610-2609f7a7-087b-4d85-b4d4-c04ede5f91e4","title":"Modernizacja kompleksu sportowego w Gminie Rossosz","organizationId":"12365","organizationName":"Gmina Rossosz","organizationPartName":null,"organizationCity":"Rossosz","organizationProvince":"lubelskie","bzpNumber":"2025/BZP 00284362/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:28:15.76Z"},{"objectId":"ocds-148610-903f87f5-7dee-42a9-8f39-bc2d0925958f","title":"Świadczenie usługi przygotowania, dostarczania i wydawania posiłków","organizationId":"4505","organizationName":"PUBLICZNA SZKOŁA PODSTAWOWA W BEZRZECZU","organizationPartName":null,"organizationCity":"Bezrzecze","organizationProvince":"zachodniopomorskie","bzpNumber":"2025/BZP 00284329/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T09:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:21:59.699Z"},{"objectId":"ocds-148610-4d98c53d-8a83-445d-927c-707ededbce7a","title":"Usługi kastracji i znakowania zwierząt właścicielskich w ramach pilotażowych programów Schroniska dla Zwierząt \"Dolina Dolistówki\" w Białymstoku","organizationId":"6558","organizationName":"Schronisko dla Zwierząt w Białymstoku","organizationPartName":null,"organizationCity":"Białystok","organizationProvince":"podlaskie","bzpNumber":"2025/BZP 00284325/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-30T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:19:53.438Z"},{"objectId":"ocds-148610-dc638293-1dc8-4462-83a0-b3ebc82ca445","title":"Dostawa aparatury lokalizacyjnej do linii kablowych SN i nn wraz z zabudową na pojazdach w latach 2025-2026 na potrzeby ENEA Operator Sp. z o. o. - 10 zadań","organizationId":"3231","organizationName":"ENEA Operator sp. z o.o.","organizationPartName":null,"organizationCity":"Poznań","organizationProvince":"wielkopolskie","bzpNumber":null,"tenderType":"2.8.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-15T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":false,"tedContractNoticeNumber":" 115/2025 392314-2025","initiationDate":"2025-06-18T08:16:33.369Z"},{"objectId":"ocds-148610-5e3b447c-a19e-4583-9663-fc395b17db87","title":"Termomodernizacja i przebudowa budynku po byłej szkole w m. Kołaczkowice wraz ze zmianą sposobu użytkowania na potrzeby klubu dziecięcego wraz z niezbędną infrastrukturą techniczną i drogową","organizationId":"4193","organizationName":"Gmina Miedźno","organizationPartName":null,"organizationCity":"Miedźno","organizationProvince":"śląskie","bzpNumber":"2025/BZP 00284310/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-07T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:16:01.112Z"},{"objectId":"ocds-148610-8788b64a-6c32-4531-bfd1-3a32e6506863","title":"Cyberbezpieczny Samorząd w Gminie Rokiciny","organizationId":"4000","organizationName":"Gmina Rokiciny","organizationPartName":null,"organizationCity":"Rokiciny","organizationProvince":"łódzkie","bzpNumber":"2025/BZP 00284294/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-26T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:13:10.919Z"},{"objectId":"ocds-148610-70072d90-b0a1-4bfb-b77a-72f354ee8232","title":"Świadczenie usługi przygotowania, dostarczania i wydawania posiłków","organizationId":"5054","organizationName":"PUBLICZNA SZKOŁA PODSTAWOWA IM. K.I. GAŁCZYŃSKIEGO W DOBREJ","organizationPartName":null,"organizationCity":"Dobra","organizationProvince":"zachodniopomorskie","bzpNumber":"2025/BZP 00284291/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T08:30:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:12:44.162Z"},{"objectId":"ocds-148610-3eb8b917-9054-4833-b25c-62025ce5143e","title":"„Świadczenie usług w zakresie dowozu uczniów do szkół podstawowych na terenie Gminy Czarnocin poprzez zakup biletów miesięcznych w roku szkolnym 2025/2026”","organizationId":"14076","organizationName":"Gmina Czarnocin","organizationPartName":null,"organizationCity":"Czarnocin","organizationProvince":"świętokrzyskie","bzpNumber":"2025/BZP 00284288/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-26T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:11:54.29Z"},{"objectId":"ocds-148610-39a1e540-59da-43bd-a36a-34e9e534f23b","title":"Dostawa rękawic diagnostycznych i chirurgicznych","organizationId":"14295","organizationName":"Samodzielny Publiczny Zakład Opieki Zdrowotnej Ministerstwa Spraw Wewętrznych i Administracji w Rzeszowie","organizationPartName":null,"organizationCity":"Rzeszów","organizationProvince":"podkarpackie","bzpNumber":"2025/BZP 00284280/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-26T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:10:47.936Z"},{"objectId":"ocds-148610-58eb9863-e17c-4e31-a8b4-08642fbe7015","title":"Przebudowa drogi wewnętrznej będącą własnością gminy Liw w miejscowości Jartypory","organizationId":"1923","organizationName":"GMINA LIW","organizationPartName":null,"organizationCity":"Węgrów","organizationProvince":"mazowieckie","bzpNumber":"2025/BZP 00284256/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:05:41.909Z"},{"objectId":"ocds-148610-7d069dda-01d2-4055-98e2-662bf8b3f16a","title":"Usługa sprzątania dla seniorów i osób z niepełnosprawnościami w ramach programu „Czysty Dom”","organizationId":"2510","organizationName":"Gmina Miejska Kraków - Urząd Miasta Krakowa","organizationPartName":null,"organizationCity":"Kraków","organizationProvince":"małopolskie","bzpNumber":"2025/BZP 00284245/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-30T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:02:55.804Z"},{"objectId":"ocds-148610-dfb69843-eb96-4f21-9bd1-bcbb593944d4","title":"Dostawa pomocy dydaktycznych na potrzeby jednostek podległych Gminie Stalowa Wola","organizationId":"8550","organizationName":"Stalowowolskie Centrum Usług Wspólnych","organizationPartName":null,"organizationCity":"Stalowa Wola","organizationProvince":"podkarpackie","bzpNumber":"2025/BZP 00284243/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-08T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:02:49.326Z"},{"objectId":"ocds-148610-c7f418a9-da01-4a87-a5e8-f323f1ec1ab0","title":"Modernizacja zjeżdżalni na pływalni WODNIK w Ozorkowie – modernizacja wanny hamownej na zjeżdżalni","organizationId":"144282","organizationName":"Centrum Sportu i Rekreacji \"Wodnik\" w Ozorkowie","organizationPartName":null,"organizationCity":"Ozorków","organizationProvince":"łódzkie","bzpNumber":"2025/BZP 00284239/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T06:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:01:43.079Z"},{"objectId":"ocds-148610-1bc4f276-1cb0-4055-b562-1421d9951f3c","title":"Modernizacja drogi gminnej nr 005504F w m. Okunin","organizationId":"2396","organizationName":"Gmina Sulechów","organizationPartName":null,"organizationCity":"Sulechów","organizationProvince":"lubuskie","bzpNumber":"2025/BZP 00284235/01","tenderType":"1.1.2","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-04T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T08:00:46.037Z"},{"objectId":"ocds-148610-72067029-c731-4332-94ed-2108fe1ab4db","title":"Doposażenie specjalistycznych pracowni dydaktycznych w Zespole Szkół Elektryczno-Elektronicznych im. M. T. Hubera w Szczecinie","organizationId":"1759","organizationName":"GMINA MIASTO SZCZECIN","organizationPartName":"Biuro ds. Zamówień Publicznych Urzędu Miasta Szczecin","organizationCity":"Szczecin","organizationProvince":"zachodniopomorskie","bzpNumber":null,"tenderType":"2.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-17T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":false,"tedContractNoticeNumber":"Dz.U. S: 115/2025 392864-2025","initiationDate":"2025-06-18T08:00:34.827Z"},{"objectId":"ocds-148610-a9a55f18-9fec-4f70-860e-c41b6f92674f","title":"Modernizacja parkingu przy budynku Urzędu Miejskiego w Ozorkowie przy ul. Listopadowej 16","organizationId":"2068","organizationName":"Gmina Miasto Ozorków","organizationPartName":null,"organizationCity":"Ozorków","organizationProvince":"łódzkie","bzpNumber":"2025/BZP 00284228/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:58:58.796Z"},{"objectId":"ocds-148610-73e377db-5d9c-4a19-91b2-00695e321f73","title":"Modernizacja kompleksu sportowego „Moje Boisko-Orlik 2012” przy ul. 19 Stycznia w Szamocinie","organizationId":"5184","organizationName":"Gmina Szamocin","organizationPartName":null,"organizationCity":"Szamocin","organizationProvince":"wielkopolskie","bzpNumber":"2025/BZP 00284225/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T09:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:58:41.746Z"},{"objectId":"ocds-148610-cd36fde0-827d-48d2-9ef5-de110645ae7e","title":"Dostawa pomocy dydaktycznych do pracowni zawodowych  Zespołu Szkół im. Walerego Goetla w Suchej Beskidzkiej","organizationId":"49","organizationName":"Powiat Suski","organizationPartName":null,"organizationCity":"Sucha Beskidzka","organizationProvince":"małopolskie","bzpNumber":"2025/BZP 00284219/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:57:35.028Z"},{"objectId":"ocds-148610-6bd5042e-22ca-40be-aa38-fa8e99c46a40","title":"Wykonanie nakładki z betonu asfaltowego na ul. Północnej w Gnieźnie","organizationId":"2885","organizationName":"Miasto Gniezno","organizationPartName":"Urząd Miejski w Gnieźnie","organizationCity":"Gniezno","organizationProvince":"wielkopolskie","bzpNumber":"2025/BZP 00284212/01","tenderType":"1.1.2","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-07T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:56:07.702Z"},{"objectId":"ocds-148610-90832b51-1e41-4d6e-a04f-c46814f4644e","title":"Sukcesywna dostawa odczynników chemicznych oraz materiałów i akcesoriów eksploatacyjnych dla Instytutu Biocybernetyki                            i Inżynierii Biomedycznej im. Macieja Nałęcza PAN","organizationId":"7510","organizationName":"Instytut Biocybernetyki i Inżynierii Biomedycznej im. Macieja Nałęcza Polskiej Akademii Nauk, ul. Księcia Trojdena 4, 02 - 109 Warszawa","organizationPartName":null,"organizationCity":"Warszawa","organizationProvince":"mazowieckie","bzpNumber":"2025/BZP 00284211/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-09T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:55:57.612Z"},{"objectId":"ocds-148610-027ca75e-e9e0-49b3-b7c6-1f7265997973","title":"Dostawa różnego sprzętu kwaterunkowego-w 8 pakietach","organizationId":"6460","organizationName":"Rejonowy Zarząd Infrastruktury w Lublinie","organizationPartName":null,"organizationCity":"Lublin","organizationProvince":"lubelskie","bzpNumber":null,"tenderType":"2.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-23T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":false,"tedContractNoticeNumber":"2025-OJS115-00395182","initiationDate":"2025-06-18T07:54:20Z"},{"objectId":"ocds-148610-ffd88281-fa72-4d6a-abab-f69f9be74789","title":"Przedłużenie wsparcia dla licencji systemu kopii zapasowych CommVault oraz przedłużenie wsparcia serwisowego biblioteki taśmowej","organizationId":"5955","organizationName":"SZKOŁA GŁÓWNA GOSPODARSTWA WIEJSKIEGO W WARSZAWIE","organizationPartName":"Sekcja Zamówień Publicznych","organizationCity":"Warszawa","organizationProvince":"mazowieckie","bzpNumber":"2025/BZP 00284188/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-01T08:30:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:51:04.734Z"},{"objectId":"ocds-148610-211196d1-41a7-47a5-863c-07145348cdf0","title":"Dowóz dzieci niepełnosprawnych do szkół oraz placówek oświatowych \nwraz z zapewnieniem opieki podczas przewozu”","organizationId":"4470","organizationName":"Gmina Miejska Pabianice reprezentowana przez Prezydenta Miasta Pabianic Grzegorza Mackiewicza","organizationPartName":null,"organizationCity":"Pabianice","organizationProvince":"łódzkie","bzpNumber":"2025/BZP 00284182/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-30T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:49:37.896Z"},{"objectId":"ocds-148610-8059fdca-4052-474c-8308-7cc21cfce34a","title":"Świadczenie usługi przygotowania, dostarczania i wydawania posiłków","organizationId":"5173","organizationName":"PUBLICZNA SZKOŁA PODSTAWOWA IM. PRZYJACIÓŁ DZIECI W MIERZYNIE","organizationPartName":null,"organizationCity":"Mierzyn","organizationProvince":"zachodniopomorskie","bzpNumber":"2025/BZP 00284167/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:48:12.076Z"},{"objectId":"ocds-148610-ca5a091a-1237-481a-bb83-13843d5b6568","title":"Dostawa artykułów żywnościowych do stołówki szkolnej w Zespole Szkół Sportowych w Siemianowicach Śląskich w roku szkolnym 2025/2026 – od września do grudnia 2025r.","organizationId":"8136","organizationName":"Zespół Szkół Sportowych w Siemianowicach Śląskich","organizationPartName":null,"organizationCity":"Siemianowice Śląskie","organizationProvince":"śląskie","bzpNumber":"2025/BZP 00284164/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-26T09:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:47:46.294Z"},{"objectId":"ocds-148610-9d0b00ec-736c-4a33-9654-5101e590bc65","title":"Dostawa i montaż wodomierzy z systemem zdalnego odczytu oraz urządzeń monitorowania stanu sieci wodociągowej","organizationId":"2242","organizationName":"Gmina Łopiennik Górny","organizationPartName":null,"organizationCity":"Łopiennik Górny","organizationProvince":"lubelskie","bzpNumber":"2025/BZP 00284159/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:46:27.603Z"},{"objectId":"ocds-148610-e2112e63-6ff8-4165-b1ef-519544034981","title":"Wykonanie kompleksowej dokumentacji projektowo-kosztorysowej budynków w Dąbrowie Górniczej w ramach zadania: Termomodernizacja gminnych budynków mieszkalnych - 3 zadania","organizationId":"3316","organizationName":"Samorządowy Zakład Budżetowy Miejski Zarząd Budynków Mieszkalnych","organizationPartName":null,"organizationCity":"Dąbrowa Górnicza","organizationProvince":"śląskie","bzpNumber":"2025/BZP 00284117/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-06-27T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:38:40.39Z"},{"objectId":"ocds-148610-1a25db84-a090-4fd4-a0d4-384dcb0601b2","title":"Mundur wyjściowy","organizationId":"11472","organizationName":"Zakład Produkcyjno Usługowo Handlowy Lasów Państwowych","organizationPartName":"ZPUH","organizationCity":"Olsztyn","organizationProvince":"warmińsko-mazurskie","bzpNumber":null,"tenderType":"2.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T07:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":false,"tedContractNoticeNumber":"OJ S 115/2025 poz.  393825-2025","initiationDate":"2025-06-18T07:37:04.007Z"},{"objectId":"ocds-148610-f2b1eade-e92c-4047-ba72-38978a835e8c","title":"Rozbudowa systemów monitoringu wizyjnego w budynkach Śląskiego Urzędu Wojewódzkiego ","organizationId":"12759","organizationName":"Śląski Urząd Wojewódzki","organizationPartName":null,"organizationCity":"Katowice","organizationProvince":"śląskie","bzpNumber":null,"tenderType":"2.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-25T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":false,"tedContractNoticeNumber":"391805-2025","initiationDate":"2025-06-18T07:36:56.325Z"},{"objectId":"ocds-148610-16d774a9-4589-46b5-9e1e-358719caa0e7","title":"Realizacja kampanii edukacyjno-informacyjnych na rzecz kształcenia zawodowego, szkolnictwa wyższego oraz uczenia się przez całe życie, w tym uczenia się dorosłych pn. „Małopolska uczy”","organizationId":"3908","organizationName":"WOJEWÓDZTWO MAŁOPOLSKIE","organizationPartName":null,"organizationCity":"Kraków","organizationProvince":"małopolskie","bzpNumber":null,"tenderType":"2.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-03T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":false,"tedContractNoticeNumber":"2025-OJS115-00392226-pl-ts","initiationDate":"2025-06-18T07:35:55.198Z"},{"objectId":"ocds-148610-6a7e34cf-a778-4050-9ff7-9627ec7c8204","title":"Przebudowa drogi gminnej nr 400565W Plac Wolności w Szydłowcu- II etap","organizationId":"5188","organizationName":"GMINA SZYDŁOWIEC","organizationPartName":null,"organizationCity":"Szydłowiec","organizationProvince":"mazowieckie","bzpNumber":"2025/BZP 00284104/01","tenderType":"1.1.1","competitionType":null,"concessionType":null,"submissionOffersDate":"2025-07-04T08:00:00Z","tenderState":"Initiated","isTenderAmountBelowEU":true,"tedContractNoticeNumber":null,"initiationDate":"2025-06-18T07:35:25.796Z"}]
type orderDTO struct {
	ObjectId                string    `json:"objectId"`
	Title                   string    `json:"title"`
	OrganizationId          string    `json:"organizationId"`
	OrganizationName        string    `json:"organizationName"`
	OrganizationPartName    string    `json:"organizationPartName"`
	OrganizationCity        string    `json:"organizationCity"`
	OrganizationProvince    string    `json:"organizationProvince"`
	BzpNumber               string    `json:"bzpNumber"`
	TenderType              string    `json:"tenderType"`
	CompetitionType         string    `json:"competitionType"`
	ConcessionType          string    `json:"concessionType"`
	SubmissionOffersDate    time.Time `json:"submissionOffersDate"`
	TenderState             string    `json:"tenderState"`
	IsTenderAmountBelowEU   bool      `json:"isTenderAmountBelowEU"`
	TedContractNoticeNumber string    `json:"tedContractNoticeNumber"`
	InitiationDate          time.Time `json:"initiationDate"`
}

func (order orderDTO) getTenderDTO() *tenderDTO {
	href := "https://ezamowienia.gov.pl/mp-client/search/list/" + order.ObjectId
	return newTenderDTO(order.Title, href, order.SubmissionOffersDate.Format("2006.01.02"), order.ObjectId)
}

func main() {
	flags := newFlagDTO()
	processTenders(flags)
	processOrders(flags)
}

func fileDateStr() string {
	return time.Now().Format("20060102")
}

func processTenders(flags *flagDTO) {
	var err error
	var done bool
	var fileOldAll *xlsx.File

	tenders := make([]*tenderDTO, 0)
	tendersIT := make([]*tenderDTO, 0)
	tendersOldAll := make([]*tenderDTO, 0)

	fileOldAll, err = xlsx.OpenFile(flags.tenderOldFileName)
	if err == nil {
		tendersOldAll = readOldAll("przetargi", fileOldAll, tendersOldAll)
	} else {
		fmt.Printf("file %s was not found\n", flags.tenderOldFileName)
		tendersOldAll = make([]*tenderDTO, 0)
		//tendersOldAll = append(tendersOldAll, newTenderDTO("przetarg", "link", "data"))
	}
	fmt.Printf("tendersOldAll len: %v\n", len(tendersOldAll))

	session := azuretls.NewSession()
	for page := 1; page <= flags.tenderPages; page++ {
		fmt.Println("tender page: ", page)
		err, tendersIT, tenders, done = processGetTenderPage(page, session, tendersIT, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	processSaveDataToExcel("przetargi", err, tenders, tendersIT, tendersOldAll, flags)
	fmt.Println("tenders END")
}

func processOrders(flags *flagDTO) {
	var err error
	var done bool
	var fileOldAll *xlsx.File

	fmt.Println("orders START")
	tenders := make([]*tenderDTO, 0)
	tendersIT := make([]*tenderDTO, 0)
	tendersOldAll := make([]*tenderDTO, 0)

	fileOldAll, err = xlsx.OpenFile(flags.ordersOldFileName)
	if err == nil {
		tendersOldAll = readOldAll("oferty", fileOldAll, tendersOldAll)
	} else {
		fmt.Printf("file %s was not found\n", flags.ordersOldFileName)
		tendersOldAll = make([]*tenderDTO, 0)
		//tendersOldAll = append(tendersOldAll, newTenderDTO("przetarg", "link", "data"))
	}
	fmt.Printf("ordersOldAll len: %v\n", len(tendersOldAll))

	session := azuretls.NewSession()
	for page := 1; page <= flags.orderPages; page++ {
		fmt.Println("order page: ", page)
		err, tendersIT, tenders, done = processGetOrderPage(page, session, tendersIT, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	processSaveDataToExcel("oferty", err, tenders, tendersIT, tendersOldAll, flags)
	fmt.Println("orders END")
}

func processSaveDataToExcel(filename string, err error, tenders, tendersIT, tendersOldAll []*tenderDTO, flags *flagDTO) {
	var fileAll *xlsx.File
	var fileIT *xlsx.File

	fmt.Println("processSaveDataToExcel")

	fileIT = xlsx.NewFile()
	err = processSaveToExcel(filename+" IT", fileIT, tendersIT, tendersOldAll)

	err = fileIT.Save(filename + "_it_" + fileDateStr() + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}

	if flags.saveAll {
		if flags.appendAll {
			for _, tender := range tendersOldAll {
				if !contains(tenders, tender) {
					tenders = append(tenders, tender)
				} else {
					fmt.Println("processSaveDataToExcel saveAll: there is already this old tender")
				}
			}
		}
		fileAll = xlsx.NewFile()
		err = processSaveAllToExcel(filename, tenders, tendersOldAll, fileAll)
		err = fileAll.Save(filename + "_all.xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}
		err = fileAll.Save(filename + "_all_" + fileDateStr() + ".xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}

func readOldAll(sheetName string, fileOldAll *xlsx.File, tendersOldAll []*tenderDTO) []*tenderDTO {
	sheet, ok := fileOldAll.Sheet[sheetName]
	if !ok {
		panic(errors.New("Sheet tenders not found"))
	}
	fmt.Println("Max row is", sheet.MaxRow)
	for row := 1; row < sheet.MaxRow; row++ {
		r, err := sheet.Row(row)
		if err != nil {
			panic(err)
		}
		tendersOldAll = oldAllRowVisitor(r, tendersOldAll)
	}
	return tendersOldAll
}

func oldAllRowVisitor(r *xlsx.Row, tendersOldAll []*tenderDTO) []*tenderDTO {
	nr := 1
	idCell := r.GetCell(nr)
	idValue := idCell.Value
	dateCell := r.GetCell(nr + 1)
	dateValue := dateCell.Value
	hrefCell := r.GetCell(nr + 2)
	hrefValue := hrefCell.Value
	nameCell := r.GetCell(nr + 3)
	nameValue := nameCell.Value
	tender := newTenderDTO(nameValue, hrefValue, dateValue, idValue)
	tendersOldAll = append(tendersOldAll, tender)
	return tendersOldAll
}

func getHrefID(value string) string {
	// _noticeId=3108196
	if len(value) < 10 {
		return "len err"
	}
	pos := strings.Index(value, "_noticeId=")
	if pos == -1 {
		return "index err"
	}
	id := value[pos+10 : len(value)]
	return id
}

func processSaveToExcel(sheetName string, file *xlsx.File, tendersIT []*tenderDTO, tendersOldAll []*tenderDTO) error {
	sheetIT, err := file.AddSheet(sheetName)
	setHeader(0, sheetIT)
	rowIT := 0
	for _, tendersT := range tendersIT {
		rowIT++
		setRowData(0, sheetIT, rowIT, tendersT)
	}
	return err
}

func processSaveAllToExcel(sheetName string, tenders, tendersOldAll []*tenderDTO, file *xlsx.File) error {
	sheet, err := file.AddSheet(sheetName)
	setAllHeader(sheet)
	rowOther := 0
	for _, tender := range tenders {
		rowOther++
		setRowData(2, sheet, rowOther, tender)
		if tender.isIT {
			cell, _ := sheet.Cell(rowOther, 0)
			cell.Value = "IT"
		}
		cell, _ := sheet.Cell(rowOther, 1)
		cell.Value = tender.id
	}
	return err
}

func setRowData(startCell int, sheet *xlsx.Sheet, r int, tender *tenderDTO) {
	nr := startCell
	cell, _ := sheet.Cell(r, nr)
	cell.Value = tender.date

	cell, _ = sheet.Cell(r, nr+1)
	cell.SetHyperlink(tender.href, tender.href, "")
	style := cell.GetStyle()
	style.Font.Underline = true
	style.Font.Color = "FF0000FF"
	cell.SetStyle(style)

	cell, _ = sheet.Cell(r, nr+2)
	cell.Value = tender.name
}

func setHeader(startCell int, sheet *xlsx.Sheet) {
	nr := startCell
	cell, _ := sheet.Cell(0, nr)
	cell.Value = "data"
	sheet.SetColWidth(nr+1, nr+1, 10)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "link"
	sheet.SetColWidth(nr+1, nr+1, 18)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "nazwa"
	sheet.SetColWidth(nr+1, nr+1, 100)
}

func setAllHeader(sheet *xlsx.Sheet) {
	setHeader(2, sheet)
	cell, _ := sheet.Cell(0, 0)
	cell.Value = "IT"
	sheet.SetColWidth(1, 1, 3)
	cell, _ = sheet.Cell(0, 1)
	cell.Value = "ID"
	sheet.SetColWidth(2, 2, 12.5)
}

func processGetTenderPage(page int, session *azuretls.Session, tendersIT []*tenderDTO, tenders []*tenderDTO, tendersOldAll []*tenderDTO) (error, []*tenderDTO, []*tenderDTO, bool) {
	pageStr := fmt.Sprintf("%d", page)

	//https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/cat?_7_WAR_organizationnoticeportlet_order=createDateDesc
	//why not this?        https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/cat?_7_WAR_organizationnoticeportlet_cur=1_order=createDateDesc
	//only this form is on www for page n: https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/cat?_7_WAR_organizationnoticeportlet_cur=1
	response, err := session.Get("https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/cat?_7_WAR_organizationnoticeportlet_cur=" + pageStr)
	if err != nil {
		panic(err)
	}
	//fmt.Println(response.String())

	element, err := gosoup.ParseAsHTML(response.String())
	if err != nil {
		// log/handle error
		fmt.Println("could not parse")
		return err, tendersIT, tenders, false
	}
	//fmt.Println("element:", element)

	containerElement := element.Find("div", gosoup.Attributes{"id": "_7_WAR_organizationnoticeportlet_selectNoticesSearchContainer"})
	if containerElement == nil {
		fmt.Println("could not find container element")
	}

	//fmt.Println("container element:", containerElement)

	subContainer := containerElement.Find("div", gosoup.Attributes{"class": "lfr-search-container-list"})
	if subContainer == nil {
		fmt.Println("could not find subContainer element")
	}
	//fmt.Println("subContainer element:", subContainer)

	group := subContainer.FindByTag("dl")
	if group == nil {
		fmt.Println("could not find group element")
	}
	//fmt.Println("group element:", group)

	expectedTag := "dd"
	expectedAttrKey := "data-qa-id"
	expectedAttrVal := "row"
	expectedElementsSize := 12
	elements := group.FindAll(expectedTag, gosoup.Attributes{expectedAttrKey: expectedAttrVal})
	if len(elements) != expectedElementsSize {
		fmt.Printf("wrong number of elements found: %q, expected number: %q", len(elements), expectedElementsSize)
	}

	for _, element := range elements {
		if element.Data != expectedTag {
			fmt.Printf("wrong element tag, expected: %q, actual: %q", expectedTag, element.Data)
		}
		attributeValue, ok := element.GetAttribute(expectedAttrKey)
		if !ok || attributeValue != expectedAttrVal {
			fmt.Printf("expected attribute: %q: %q does not exist", expectedAttrKey, expectedAttrVal)
		}

		//TODO add app parameter --debug
		if false {
			fmt.Println("dd element:", element)
		}

		aTag := element.FindByTag("a")
		if aTag == nil {
			fmt.Println("could not find aTag element")
		}
		hrefValue, ok := aTag.GetAttribute("href")
		if !ok {
			fmt.Printf("href attribute: does not exist")
		}

		//a tag have a name content
		nameValue := aTag.FirstChild.Data
		//class="notice-date"
		//dateDiv := element.Find("div", gosoup.Attributes{"class": "notice-date"})
		dateSpan := element.Find("span", gosoup.Attributes{"title": "Termin składania ofert"})
		dateTimeValue := strings.TrimSpace(dateSpan.FirstChild.Data)

		//t, err := time.Parse(time.RFC3339, "2023-05-02T09:34:01Z")
		//Mon Jun 23 09:00:00 GMT 2025: example value
		//Mon Jan _2 15:04:05 GMT 2006: layout form
		const longForm = "Mon Jan _2 15:04:05 GMT 2006"
		dateTime, _ := time.Parse(longForm, dateTimeValue)
		dateValue := dateTime.Format("2006.01.02")

		tender := newTenderDTO(nameValue, hrefValue, dateValue, getHrefID(hrefValue))

		tendersIT, tenders = appendTender(tender, tendersIT, tenders)
		if contains(tendersOldAll, tender) {
			fmt.Println("processGetTenderPage: old tenders contains this", tender)
			return err, tendersIT, tenders, true
		}
	}
	return err, tendersIT, tenders, false
}

func appendTender(tender *tenderDTO, tendersIT, tenders []*tenderDTO) ([]*tenderDTO, []*tenderDTO) {
	if tender.isIT {
		tendersIT = append(tendersIT, tender)
		tenders = append(tenders, tender)
	} else {
		tenders = append(tenders, tender)
	}
	return tendersIT, tenders
}

func processGetOrderPage(page int, session *azuretls.Session, tendersIT []*tenderDTO, tenders []*tenderDTO, tendersOldAll []*tenderDTO) (error, []*tenderDTO, []*tenderDTO, bool) {
	var orders []orderDTO
	pageStr := fmt.Sprintf("%d", page)
	response, err := session.Get("https://ezamowienia.gov.pl/mp-readmodels/api/Search/SearchTenders?SortingColumnName=InitiationDate&SortingDirection=DESC&PageNumber=" + pageStr + "&PageSize=50")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(response.String()), &orders)
	if err != nil {
		println(err.Error())
		println("response:" + response.String())
		return err, nil, nil, false
	}
	//fmt.Printf("order: %+v \n", orders)
	//fmt.Println(len(orders))
	for _, order := range orders {
		tender := order.getTenderDTO()
		tendersIT, tenders = appendTender(tender, tendersIT, tenders)
		if contains(tendersOldAll, tender) {
			fmt.Println("processGetOrderPage: old orders contains this:", tender)
			return err, tendersIT, tenders, true
		}
	}
	return err, tendersIT, tenders, false
}

func contains(tenders []*tenderDTO, tender *tenderDTO) bool {
	for _, p := range tenders {
		if p.id == tender.id && p.name == tender.name {
			fmt.Printf("id=%s;", p.id)
			return true
		}
	}
	return false
}
