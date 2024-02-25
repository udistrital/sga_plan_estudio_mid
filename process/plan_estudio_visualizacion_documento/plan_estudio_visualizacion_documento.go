package plan_estudio_visualizacion_documento

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/sga_plan_estudio_mid/helpers"
	"github.com/udistrital/sga_plan_estudio_mid/utils"
	"sort"
)

func GenerateStudyPlanDocument(data map[string]interface{}) *gofpdf.Fpdf {
	// page features
	pdf := gofpdf.New("L", "mm", "Legal", "")

	marginTB := 4.0
	marginLR := 4.0
	pdf.SetMargins(marginLR, marginTB, marginLR)
	pdf.SetAutoPageBreak(true, 1.0)

	pageStyle := getPageStyle(pdf)
	planMetadata := getPlanMetadata(data, pageStyle)

	// create pages
	for i := 0; i < planMetadata.numPages; i++ {
		pdf.AddPage()
	}

	for nPage := 1; nPage <= planMetadata.numPages; nPage++ {
		pdf.SetPage(nPage)
		pdf.SetHomeXY()

		// draw page margin
		x, y := pdf.GetXY()
		pdf.SetDrawColor(
			pageStyle.ComplementaryColorRGB[0],
			pageStyle.ComplementaryColorRGB[1],
			pageStyle.ComplementaryColorRGB[2])
		pdf.RoundedRect(x-1, y-1, pageStyle.WW+1, pageStyle.HW+1, 2, "1234", "D")
		pdf.SetDrawColor(0, 0, 0)

		pdf.SetXY(x, y)
		pdf = studyPlanHeader(pdf, data, pageStyle)

		// Add cards, card by project
		plans, plansOk := data["Planes"]
		if plansOk {
			x, y = pageStyle.ML+planMetadata.externalCardSpace, pageStyle.HH+3
			pdf.SetXY(x, y)
			widthCard := 0.0
			sumWidthCard := 0.0
			for nProject := 0; nProject < planMetadata.numProjects; nProject++ {
				dataProject, dataProjectOk := plans.([]any)[nProject].(map[string]any)
				if dataProjectOk {
					widthCard = planMetadata.cardStyleProject[nProject].mainCardWidth
					if planMetadata.distributionConfig.splitHorizontal && nPage != 1 {
						widthCard = planMetadata.cardStyleProject[nProject].secondaryCardWidth
						if widthCard > 0.0 {
							pdf = createProjectCard(pdf, dataProject, pageStyle, planMetadata, widthCard, nPage, nProject)
							sumWidthCard = sumWidthCard + widthCard + planMetadata.externalCardSpace
							pdf.SetXY(x+sumWidthCard, y)
						}
					} else {
						pdf = createProjectCard(pdf, dataProject, pageStyle, planMetadata, widthCard, nPage, nProject)
						sumWidthCard = sumWidthCard + widthCard + planMetadata.externalCardSpace
						pdf.SetXY(x+sumWidthCard, y)
					}
				}

			}
		}

		// Add footer
		pdf = studyPlanFooter(pdf, data, pageStyle)
	}

	return pdf
}

func getPageStyle(pdf *gofpdf.Fpdf) utils.PageStyle {
	widthPage, heightPage := pdf.GetPageSize()
	l, t, r, b := pdf.GetMargins()
	pageStyle := utils.PageStyle{
		ML: l,
		MT: t,
		MR: r,
		MB: b,
		WW: widthPage - l - r,
		HW: heightPage - (2 * t),
		HH: 30,
		HB: heightPage - 30,
		HF: 0}

	// blue for headers
	pageStyle.BaseColorRGB[0] = 20
	pageStyle.BaseColorRGB[1] = 103
	pageStyle.BaseColorRGB[2] = 143

	// gray for headers
	pageStyle.SecondaryColorRGB[0] = 128
	pageStyle.SecondaryColorRGB[1] = 128
	pageStyle.SecondaryColorRGB[2] = 128

	// light blue for outlines
	pageStyle.ComplementaryColorRGB[0] = 90
	pageStyle.ComplementaryColorRGB[1] = 149
	pageStyle.ComplementaryColorRGB[2] = 184

	return pageStyle
}

func studyPlanHeader(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Added university logo
	path := beego.AppConfig.String("StaticPath")
	x := ((pageStyle.WW - 70) / 4) - 10
	y := pageStyle.MT
	pdf = utils.AddImage(pdf, path+"/img/logoud.jpeg", x, y, 0, 20)

	// Added title and subtitle
	facultyName, facNameOk := data["Facultad"]
	if facNameOk == false || facultyName == nil {
		facultyName = ""
	}
	facultyNameSize := float64(len(fmt.Sprintf("%v", facultyName)))

	utils.FontStyle(pdf, "B", 9, 0, "Helvetica")
	pdf.SetXY(((pageStyle.WW+35-facultyNameSize)/4)-20, y+16)
	pdf.Cell(5, 10, tr(fmt.Sprintf("%v", facultyName)))
	pdf.Ln(5)

	planName, planNameOk := data["Nombre"]
	if planNameOk == false || planName == nil {
		planName = ""
	}
	planNameSize := float64(len(fmt.Sprintf("%v", planName)))

	utils.FontStyle(pdf, "", 8, 0, "Helvetica")
	y = pdf.GetY() - 1
	pdf.SetXY(((pageStyle.WW+35-planNameSize)/4)-10, y)
	pdf.Cell(5, 10, tr(fmt.Sprintf("%v", planName)))

	// Added space detail
	pathDesc := beego.AppConfig.String("StaticPath")
	x = ((pageStyle.WW - 70) / 2) + 20
	y = pageStyle.MT + 2
	pdf = utils.AddImage(pdf, pathDesc+"/img/space_academic_detail_footer_es.png", x, y, 0, 24)

	// Added academy colors
	x = pdf.GetX() + x + 4.5
	r, g, b := pdf.GetDrawColor()
	squareWidth := 80.0
	squareHeight := 22.0
	pdf.SetDrawColor(
		pageStyle.SecondaryColorRGB[0],
		pageStyle.SecondaryColorRGB[1],
		pageStyle.SecondaryColorRGB[2])
	pdf.RoundedRect(x, y+0.5, squareWidth, squareHeight, 0, "1234", "D")

	academies, academiesOk := data["Escuelas"]
	if academiesOk != false || fmt.Sprintf("%v", academies) != "map[]" || academies != nil {
		// sort academies by quantity
		keysAcademies := make([]string, 0, len(academies.(map[string]any)))
		var academiesSort []map[string]any

		for k := range academies.(map[string]any) {
			keysAcademies = append(keysAcademies, k)
		}
		sort.SliceStable(keysAcademies, func(i, j int) bool {
			return academies.(map[string]any)[keysAcademies[i]].(map[string]any)["Cantidad"].(float64) > academies.(map[string]any)[keysAcademies[j]].(map[string]any)["Cantidad"].(float64)
		})

		for _, academy := range keysAcademies {
			academiesSort = append(academiesSort, academies.(map[string]any)[academy].(map[string]any))
		}

		academyWidth := (squareWidth / 2) - 2
		academyHeight := (squareHeight - 6) / 5
		rFont, gFont, bFont := pdf.GetTextColor()
		x = x + 1.5
		yInitialAcademy := y + 1.5
		y = yInitialAcademy
		yIncrement := 0.0
		for iAcademy, element := range academiesSort {
			utils.FontStyle(pdf, "B", 6.5, 255, "Helvetica")
			colorRGB, err := utils.Hex2RGB(fmt.Sprintf("%v", element["Color"]))
			if err != nil {
				colorRGB = [3]int{255, 255, 255}
			}

			colorFontRGB, errFont := utils.Hex2RGB(fmt.Sprintf("%v", element["TxtColor"]))
			if errFont != nil {
				colorFontRGB = [3]int{0, 0, 0}
			}

			pdf.SetFillColor(
				colorRGB[0],
				colorRGB[1],
				colorRGB[2])
			pdf.SetTextColor(colorFontRGB[0], colorFontRGB[1], colorFontRGB[2])
			if iAcademy > 4 {
				x = x + academyWidth + 1
				y = yInitialAcademy
				yIncrement = 0.0
			}

			pdf.SetXY(x, y+((academyHeight+1.0)*yIncrement))
			pdf.CellFormat(
				academyWidth, academyHeight,
				tr(fmt.Sprintf("%v", helpers.DefaultToMapString(element, "Nombre", ""))),
				"1", 1, "CM", true, 0, "")
			yIncrement++
		}
		pdf.SetTextColor(rFont, gFont, bFont)
	}
	pdf.SetDrawColor(r, g, b)
	return pdf
}

func createProjectCard(pdf *gofpdf.Fpdf, dataProject map[string]interface{}, pageStyle utils.PageStyle, planMetadata PlanMetadata, widthCard float64, nPage, nProject int) *gofpdf.Fpdf {
	x, y := pdf.GetXY()
	initX := x

	// draw card margin
	pdf.SetDrawColor(
		pageStyle.SecondaryColorRGB[0],
		pageStyle.SecondaryColorRGB[1],
		pageStyle.SecondaryColorRGB[2])
	pdf.RoundedRect(x, y-1, widthCard, pageStyle.HB-7.5, 2, "1234", "D")

	x = x + widthCard/2.0 - 22.0
	if planMetadata.doubleCol {
		pdf.SetXY(x, y+0.7)
	} else {
		pdf.SetXY(x, y+0.4)
	}
	pdf = createProjectInformationTable(pdf, dataProject, pageStyle, planMetadata.doubleCol)

	pdf.SetX(initX)
	pdf = createProjectDetails(pdf, dataProject, pageStyle, planMetadata, widthCard, nPage, nProject)

	pdf.SetX(x + 1)
	pdf = createTotalProjectCreditTable(pdf, dataProject, pageStyle, planMetadata.doubleCol)
	return pdf
}

// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
// FUNCIONES PARA CREAR TARJETA CON
// EL CONTENIDO DE CADA PROYECTO
// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

func createProjectInformationTable(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, doubleCol bool) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	infLabels := map[string]interface{}{
		"es": []string{
			tr("Resolución de aprobación"),
			tr("Total Créditos"),
			tr("Código SNIES"),
			tr("Plan de estudios")},
		"en": []string{
			tr("Approval resolution"),
			tr("Total Credits"),
			tr("SNIES code"),
			tr("Study plan")},
	}

	x := pdf.GetX()
	cellWidth := 46.0
	doubleColBorder := "B"
	if doubleCol {
		cellWidth = 92
		x = x - 24
		pdf.SetX(x)
		doubleColBorder = "BR"
	}
	cellHeight := 3.0

	// Header
	utils.FontStyle(pdf, "B", 6.5, 255, "Helvetica")
	pdf.SetFillColor(
		pageStyle.BaseColorRGB[0],
		pageStyle.BaseColorRGB[1],
		pageStyle.BaseColorRGB[2])
	pdf.CellFormat(
		cellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Nombre", ""))),
		"", 1, "CM", true, 0, "")

	// Body
	pdf.SetFillColor(0, 0, 0)
	utils.FontStyle(pdf, "", 6, 0, "Helvetica")

	pdf.SetX(x)
	bodyCellWidth := float64(int(cellWidth * 0.65))
	if doubleCol {
		cellWidth = float64(int(cellWidth / 2))
		bodyCellWidth = float64(int(cellWidth * 0.65))
	}

	pdf.CellFormat(bodyCellWidth, cellHeight, fmt.Sprintf("%v", infLabels["es"].([]string)[0]), "B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Resolucion", ""))),
		doubleColBorder, 0, "LM", false, 0, "")

	if !doubleCol {
		pdf.Ln(-1)
		pdf.SetX(x)
	}

	pdf.CellFormat(bodyCellWidth, cellHeight, fmt.Sprintf("%v", infLabels["es"].([]string)[1]), "B", 0, "LM", false, 0, "")
	pdf.CellFormat(cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Creditos", 0.0))),
		"B", 1, "LM", false, 0, "")

	pdf.SetX(x)
	pdf.CellFormat(bodyCellWidth, cellHeight, fmt.Sprintf("%v", infLabels["es"].([]string)[2]), "B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Snies", ""))),
		doubleColBorder, 0, "LM", false, 0, "")

	if !doubleCol {
		pdf.Ln(-1)
		pdf.SetX(x)
	}

	pdf.CellFormat(bodyCellWidth, cellHeight, fmt.Sprintf("%v", infLabels["es"].([]string)[3]), "B", 0, "LM", false, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "PlanEstudio", ""))),
		"B", 1, "LM", false, 0, "")
	y := pdf.GetY()
	pdf.SetXY(x, y-1.5)
	return pdf
}

func createProjectDetails(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, planMetadata PlanMetadata, widthCard float64, nPage, nProject int) *gofpdf.Fpdf {
	height := 142.5

	x, y := pdf.GetXY()
	initialPointX, initialPointY := x, y

	// draw project margin
	x = x + 2
	pdf.SetX(x)
	pdf.SetDrawColor(
		pageStyle.ComplementaryColorRGB[0],
		pageStyle.ComplementaryColorRGB[1],
		pageStyle.ComplementaryColorRGB[2])
	pdf.RoundedRect(x, y+3, widthCard-4, height, 2, "1234", "D")

	totalPeriods := planMetadata.cardStyleProject[nProject].periodsMainPage
	initialPeriod := 0
	if planMetadata.distributionConfig.splitHorizontal {
		if nPage != 1 {
			totalPeriods = planMetadata.cardStyleProject[nProject].periodsSecondaryPage
			initialPeriod = planMetadata.cardStyleProject[nProject].periodsMainPage
		}
	}
	infoPeriods, infoPeriodsOk := data["InfoPeriodos"]

	if infoPeriodsOk {
		currentPeriods := infoPeriods.([]any)[initialPeriod : initialPeriod+totalPeriods]

		nPeriod := 0
		for numPer := 0; numPer < totalPeriods; numPer++ {
			nPeriod = planMetadata.cardStyleProject[nProject].initialPeriodNum + initialPeriod + numPer
			x = initialPointX + (float64(numPer) * (planMetadata.distributionConfig.colWidth + planMetadata.distributionConfig.colSpacing)) + planMetadata.distributionConfig.outerSpace - 1

			pdf.SetXY(x, initialPointY)
			pdf = createPeriod(pdf, currentPeriods[numPer].(map[string]any), pageStyle, planMetadata, widthCard, nPeriod, nPage)
		}
	}
	pdf.SetXY(initialPointX, initialPointY+height+4)

	return pdf
}

func createTotalProjectCreditTable(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, doubleCol bool) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	infLabels := map[string]interface{}{
		"es": []string{
			tr("Ítem"),
			tr("Total"),
			tr("Obligatorio Básico"),
			tr("Obligatorio Complementario"),
			tr("Electiva Intrínseca"),
			tr("Electiva Extrínseca")},
		"en": []string{
			tr("Item"),
			tr("Total"),
			tr("Basic Required"),
			tr("Complementary Required"),
			tr("Intrinsic Elective"),
			tr("Extrinsic Elective")},
	}
	x := pdf.GetX()

	if doubleCol {
		x = x - 20
		pdf.SetX(x)
	}
	cellWidth := 40.0
	cellHeight := 3.0
	bodyCellWidth := float64(int(cellWidth * 0.75))

	// Header
	utils.FontStyle(pdf, "B", 6, 255, "Helvetica")
	pdf.SetFillColor(
		pageStyle.SecondaryColorRGB[0],
		pageStyle.SecondaryColorRGB[1],
		pageStyle.SecondaryColorRGB[2])
	pdf.CellFormat(
		bodyCellWidth, cellHeight+0.5,
		fmt.Sprintf("%v", infLabels["es"].([]string)[0]),
		"", 0, "CM", true, 0, "")
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight+0.5,
		fmt.Sprintf("%v", infLabels["es"].([]string)[1]),
		"", 0, "CM", true, 0, "")

	if doubleCol {
		pdf.CellFormat(
			bodyCellWidth, cellHeight+0.5,
			fmt.Sprintf("%v", infLabels["es"].([]string)[0]),
			"", 0, "CM", true, 0, "")
		pdf.CellFormat(
			cellWidth-bodyCellWidth, cellHeight+0.5,
			fmt.Sprintf("%v", infLabels["es"].([]string)[1]),
			"", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)

	// Body
	summaryData, summaryDataOk := data["Resumen"]
	obValue := 0
	ocValue := 0
	eiValue := 0
	eeValue := 0

	if summaryDataOk {
		val, valOk := summaryData.(map[string]any)["OB"]
		if valOk {
			obValue = int(val.(float64))
		}

		val, valOk = summaryData.(map[string]any)["OC"]
		if valOk {
			ocValue = int(val.(float64))
		}

		val, valOk = summaryData.(map[string]any)["EI"]
		if valOk {
			eiValue = int(val.(float64))
		}

		val, valOk = summaryData.(map[string]any)["EE"]
		if valOk {
			eeValue = int(val.(float64))
		}
	}

	pdf.SetFillColor(0, 0, 0)
	utils.FontStyle(pdf, "", 6, 0, "Helvetica")

	pdf.SetX(x)
	// Label OB
	pdf.CellFormat(
		bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", infLabels["es"].([]string)[2]),
		"B", 0, "LM", false, 0, "")
	// Value
	pdf.CellFormat(cellWidth-bodyCellWidth, cellHeight, fmt.Sprintf("%v", obValue),
		"B", 0, "CM", false, 0, "")

	if !doubleCol {
		pdf.Ln(-1)
		pdf.SetX(x)
	}
	// Label OC
	pdf.CellFormat(
		bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", infLabels["es"].([]string)[3]),
		"B", 0, "LM", false, 0, "")
	// Value
	pdf.CellFormat(cellWidth-bodyCellWidth, cellHeight, fmt.Sprintf("%v", ocValue),
		"B", 1, "CM", false, 0, "")

	pdf.SetX(x)
	// Label EI
	pdf.CellFormat(
		bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", infLabels["es"].([]string)[4]),
		"B", 0, "LM", false, 0, "")
	// Value
	pdf.CellFormat(cellWidth-bodyCellWidth, cellHeight, fmt.Sprintf("%v", eiValue),
		"B", 0, "CM", false, 0, "")

	if !doubleCol {
		pdf.Ln(-1)
		pdf.SetX(x)
	}
	// Label EE
	pdf.CellFormat(
		bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", infLabels["es"].([]string)[5]),
		"B", 0, "LM", false, 0, "")
	// Value
	pdf.CellFormat(
		cellWidth-bodyCellWidth, cellHeight,
		fmt.Sprintf("%v", eeValue),
		"B", 1, "CM", false, 0, "")
	return pdf
}

// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
// FUNCIONES PARA CREAR TARJETA CON
// EL CONTENIDO DE CADA PROYECTO
// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

func createAcademicSpaceTable(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, planMetadata PlanMetadata) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	var colorSpace [3]int
	var colorText [3]int

	h2 := planMetadata.distributionConfig.rowHeight
	h1 := h2 * 2

	tableWidth := planMetadata.distributionConfig.colWidth
	codWidth := float64(int(tableWidth * 0.3))
	w2 := tableWidth / 5.0
	initialPointX := pdf.GetX()

	colorRGBBackground, err := utils.Hex2RGB(
		fmt.Sprintf("%v", helpers.DefaultToMapString(
			data, "EscuelaColor", "#FFFFFF")))
	if err != nil {
		colorSpace = [3]int{255, 255, 255}
	} else {
		colorSpace = colorRGBBackground
	}

	colorRGBText, err := utils.Hex2RGB(
		fmt.Sprintf("%v", helpers.DefaultToMapString(
			data, "EscuelaColorFuente", "#000000")))
	if err != nil {
		colorText = [3]int{0, 0, 0}
	} else {
		colorText = colorRGBText
	}

	pdf.SetFillColor(
		colorSpace[0],
		colorSpace[1],
		colorSpace[2])
	pdf.SetTextColor(colorText[0], colorText[1], colorText[2])
	pdf.CellFormat(
		codWidth, h1,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Codigo", ""))),
		"LT", 0, "CM", true, 0, "")

	// Celda del nombre espacio académico
	spaceName := helpers.DefaultToMapString(data, "Nombre", "")
	spaceNameList := pdf.SplitLines(
		[]byte(fmt.Sprintf("%v", spaceName)),
		tableWidth-codWidth-2)

	var borderStr string
	h11 := h1 / float64(len(spaceNameList))
	x := pdf.GetX()
	for i := 0; i < len(spaceNameList); i++ {
		pdf.SetX(x)
		if i == 0 {
			borderStr = "LTR"
		} else {
			borderStr = "LR"
		}

		pdf.CellFormat(
			tableWidth-codWidth, h11,
			tr(fmt.Sprintf("%v", string(spaceNameList[i]))),
			borderStr, 1, "CM", true, 0, "")
	}

	// celdas horas de trabajo y clasificación espacio académico
	pdf.SetX(initialPointX)

	credits, creditsOk := data["Creditos"]
	if !creditsOk {
		credits = 0
	}
	pdf.CellFormat(w2, h2, fmt.Sprintf("%v", credits),
		"LTR", 0, "CM", true, 0, "")

	htd, htdOk := data["HTD"]
	if !htdOk {
		htd = 0
	}
	pdf.CellFormat(w2, h2, fmt.Sprintf("%v", htd),
		"LTR", 0, "CM", true, 0, "")

	htc, htcOk := data["HTC"]
	if !htcOk {
		htc = 0
	}
	pdf.CellFormat(w2, h2, fmt.Sprintf("%v", htc),
		"LTR", 0, "CM", true, 0, "")

	hta, htaOk := data["HTA"]
	if !htaOk {
		hta = 0
	}
	pdf.CellFormat(w2, h2, fmt.Sprintf("%v", hta),
		"LTR", 0, "CM", true, 0, "")

	pdf.CellFormat(
		w2, h2,
		tr(fmt.Sprintf("%v", helpers.DefaultToMapString(data, "Clasificacion", ""))),
		"LTR", 1, "CM", true, 0, "")

	// Celda prerequisitos
	prerequisites, prerequisitesOk := data["Prerequisitos"]
	prerequisitesStr := ""
	if prerequisitesOk && len(prerequisites.([]any)) > 0 {
		for ipr, preRQ := range prerequisites.([]any) {
			if ipr == 0 {
				prerequisitesStr = preRQ.(string)
			} else {
				prerequisitesStr = fmt.Sprintf("%v, %v", prerequisitesStr, preRQ)
			}
		}
	}

	prerequisitesList := pdf.SplitLines(
		[]byte(fmt.Sprintf("%v", prerequisitesStr)),
		tableWidth-2)
	if len(prerequisitesList) == 0 {
		borderStr = "LTRB"
		pdf.SetX(initialPointX)
		pdf.CellFormat(
			tableWidth, h2, "",
			borderStr, 1, "CM", true, 0, "")
	} else {
		h21 := h2 / float64(len(prerequisitesList))
		for i := 0; i < len(prerequisitesList); i++ {
			pdf.SetX(initialPointX)
			if i == 0 {
				borderStr = "LTR"
			} else {
				borderStr = "LR"
			}

			if i == len(prerequisitesList)-1 {
				borderStr = fmt.Sprintf("%vB", borderStr)
			}
			pdf.CellFormat(
				tableWidth, h21,
				tr(fmt.Sprintf("%v", string(prerequisitesList[i]))),
				borderStr, 1, "CM", true, 0, "")
		}
	}

	return pdf
}

func createPeriod(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle, planMetadata PlanMetadata, widthCard float64, numPeriod, nPage int) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	x, y := pdf.GetXY()
	y = y + 5
	// Título periodo
	pdf.SetXY(x+1, y)
	utils.FontStyle(pdf, "", planMetadata.distributionConfig.fontSize, 0, "Helvetica")

	pdf.SetDrawColor(
		93,
		177,
		100)
	pdf.SetFillColor(
		255,
		255,
		255)
	pdf.CellFormat(
		planMetadata.distributionConfig.colWidth, planMetadata.distributionConfig.rowHeight,
		tr(fmt.Sprintf("Periodo: %v", numPeriod)),
		"1", 1, "CM", true, 0, "")

	// Body
	pdf.SetXY(x+1, y+planMetadata.distributionConfig.rowHeight+planMetadata.distributionConfig.rowSpacing)
	y = pdf.GetY()
	pdf.SetDrawColor(
		pageStyle.ComplementaryColorRGB[0],
		pageStyle.ComplementaryColorRGB[1],
		pageStyle.ComplementaryColorRGB[2])
	utils.FontStyle(pdf, "", planMetadata.distributionConfig.fontSize, 255, "Helvetica")

	spacesData, spacesDataOk := data["Espacios"]

	if spacesDataOk {
		totalSpaces := len(spacesData.([]any))
		numCurrentSpaces := totalSpaces
		initialNumSpace := 0
		drawSpaces := true

		if planMetadata.distributionConfig.splitVertical && totalSpaces > planMetadata.distributionConfig.maxColsRows {
			if nPage == 1 {
				numCurrentSpaces = planMetadata.distributionConfig.maxColsRows
			} else {
				initialNumSpace = planMetadata.distributionConfig.maxColsRows
				numCurrentSpaces = totalSpaces - planMetadata.distributionConfig.maxColsRows
			}
		}

		if planMetadata.distributionConfig.splitVertical && totalSpaces <= planMetadata.distributionConfig.maxColsRows && nPage != 1 {
			drawSpaces = false
		}

		if drawSpaces {
			rowHeight := (planMetadata.distributionConfig.rowHeight * 4) + planMetadata.distributionConfig.rowSpacing
			for numSpace := 0; numSpace < numCurrentSpaces; numSpace++ {
				currentSpaceData := spacesData.([]any)[initialNumSpace+numSpace]
				pdf.SetXY(x+1, y+(float64(numSpace)*rowHeight))
				pdf = createAcademicSpaceTable(pdf, currentSpaceData.(map[string]any), pageStyle, planMetadata)
			}
		}

		// Título Cantidad de créditos
		pdf.SetXY(x+1, pdf.GetY()+2)
		utils.FontStyle(pdf, "", planMetadata.distributionConfig.fontSize, 0, "Helvetica")
		pdf.SetDrawColor(
			175,
			127,
			93)
		pdf.SetFillColor(
			255,
			255,
			255)

		totalCredits, creditsOk := data["TotalCreditos"]
		if !creditsOk {
			totalCredits = 0
		}
		pdf.CellFormat(
			planMetadata.distributionConfig.colWidth, planMetadata.distributionConfig.rowHeight,
			tr(fmt.Sprintf("Cantidad de créditos: %v", totalCredits)),
			"1", 1, "CM", true, 0, "")
	} else {
		// Título Cantidad de créditos
		pdf.SetXY(x+1, pdf.GetY()+2)
		utils.FontStyle(pdf, "", planMetadata.distributionConfig.fontSize, 0, "Helvetica")
		pdf.SetDrawColor(
			175,
			127,
			93)
		pdf.SetFillColor(
			255,
			255,
			255)
		pdf.CellFormat(
			planMetadata.distributionConfig.colWidth, planMetadata.distributionConfig.rowHeight,
			tr(fmt.Sprintf("Cantidad de créditos: %v", 0)),
			"1", 1, "CM", true, 0, "")
	}
	return pdf
}

func studyPlanFooter(pdf *gofpdf.Fpdf, data map[string]interface{}, pageStyle utils.PageStyle) *gofpdf.Fpdf {
	utils.FontStyle(pdf, "", 8, 0, "Helvetica")
	pdf.SetXY(pageStyle.WW-8, pageStyle.HW+4)
	pdf.Cell(5, 3, fmt.Sprintf("PAG %v/%v", pdf.PageNo(), pdf.PageCount()))
	return pdf
}
