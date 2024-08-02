package main

import (
	"fmt"
	"html/template"

	log "github.com/sirupsen/logrus"
)

type BlogPara struct {
	ImgPath string        // img alongside
	Txt     template.HTML // text in the paragraph
}

type BlogNav struct {
	Next     string // link redirect to the next blog
	Prev     string // link redirect to the previous blog
	BaseHref string // base resource url // example /dear-diary
}

type Blog struct {
	Header   string
	Slog     string // this is when mapping to the url, or searching by the url params
	CoverImg string
	Purport  template.HTML // this is an phrase to summarise the blog but with words not directly related to the blog
	Gist     template.HTML // this gives a shrt description on what the blog is all about
	Author   string
	Owner    string
	Paras    []BlogPara
	Nav      *BlogNav // redirect away from this blog to another relevant blogs
}

type ListOfBlogs []Blog

func (lb ListOfBlogs) SearchWith(slog string) (*Blog, error) {
	for _, b := range lb {
		if b.Slog == slog {
			return &b, nil
		}
	}
	return nil, fmt.Errorf("failed to get any blog with the slog %s", slog)
}

// Paginate : for a  given page this can send back only those blogs which are relevant to the page.
// Its also sends the page information - index of the page , href url of the page
// Incase of error the result is nil, with custom error which indicates the exact error and http error code
//
/*

	result := DiaryData.Paginate(ENLISTING_PER_PAGE, page)
	if yes, herr := result.HasError(); yes {
		herr.ToHttpCtx(c)
		return
	}

	c.HTML(http.StatusOK, "page.html", gin.H{
		"Data": result.Result.(*PaginationResult),
	})
*/
func (lb ListOfBlogs) Paginate(perPage, currPage int) *ResultOrErr {
	// Input validation
	if len(lb) == 0 || perPage <= 0 {
		log.WithFields(log.Fields{
			"bloglist_len": len(lb),
			"per_page":     perPage,
		}).Error("failed pagination")
		return &ResultOrErr{Err: &InvalidArgument{Err: fmt.Errorf("bloglist empty or the number blogs per page is invalid")}, Result: nil}
	}
	// -----------
	// cCounting the total number of pages
	// -----------
	count := len(lb) / perPage
	if len(lb)%perPage != 0 {
		count++
	}
	if currPage < 1 {
		log.Warnf("requested page: %d, invalid switching to page 1", currPage)
		currPage = 1 // currPage has to be atleast 1, any less and the start will be miscaculated
	}
	// -----------------
	// getting the paged results from the entire data
	// -----------------
	start := (currPage - 1) * perPage // slice index starts from 0, pages start from 1
	end := start + perPage

	pages := []Page{}
	for i := 1; i <= count; i++ {
		// For the array of page buttons Idx on the button title , while Href used for button.onclick event
		pages = append(pages, Page{IsCurr: i == currPage, Idx: i, HRef: fmt.Sprintf("%s/?page=%d", lb[0].Nav.BaseHref, i)}) // appending page objects
	}
	// []Page -- > PaginationResult{} -- > ResultOrErr
	if currPage < count {
		return &ResultOrErr{Err: nil, Result: &PaginationResult{BlogList: lb[start:end], TotalPages: pages}}
	} else if currPage == count {
		// for the last page, it cannot have an end, since then it would be out of bounds exception
		return &ResultOrErr{Err: nil, Result: &PaginationResult{BlogList: lb[start:], TotalPages: pages}} // for the last page
	} else {
		// When its beyond the page limits then it should be an empty array
		log.WithFields(log.Fields{
			"page":        currPage,
			"pages_total": count,
		}).Errorf("page beyond number of count")
		return &ResultOrErr{Err: &InvalidQueryParam{
			Err: fmt.Errorf("invalid page number for the blog list"),
		}, Result: nil}
	}
}

var (
	DiaryData = ListOfBlogs{
		{
			Header:   "March 2024",
			Slog:     "march-2024",
			CoverImg: "tomato_gravel2.png",
			Purport:  "Spirits are high, energies are focused, and everyone is excited to begin. I have a feeling this will be a long journey, but the team is optimistic about the results. We're ready to hit the ground running.",
			Gist:     `All that happened in March of 2024`,
			Author:   "Niranjan Awati",
			Owner:    "Eensymachines, Pune",
			Nav:      &BlogNav{Next: "april-2024", Prev: "", BaseHref: "/dear-diary"},
			Paras: []BlogPara{
				{ImgPath: "", Txt: `For our initial setup, we implemented a <span class="text-dark-emphasis">1:2 ratio of bed volume to fish
				water</span>, which deviates from standard recommendations. However, we chose this approach to maintain a <span class="text-dark-emphasis">low stocking density.</span>`},
				{ImgPath: "aquaponics_bell_siphon.png", Txt: `We faced several challenges with siphon calculations, but after overcoming a few setbacks, we achieved overall stability. Perfecting the siphon demanded a significant amount of our time and manpower, as
				we were aware that once commissioned, it would be difficult to replace or even open for maintenance. At one
				stage, we had to remove a functioning siphon due to its inability to break the flow at high inlet rates.`},
				{ImgPath: "", Txt: `The water pump we selected was an old, non-submersible 1 HP motor. While it performed well, it was notably noisy and excessive for the <span class="text-dark-emphasis">48-inch head and 220-liter grow bed.</span> The
				flooding cycles were infrequent, but we allowed this over-design to persist as we were eager to see the results—optimization was reserved for a later stage.`},
				{ImgPath: "", Txt: `Next, we conducted fish-less cycling. Some might consider this overly cautious, but we wanted to avoid
				risking the fish and had the time to ensure the setup was functioning correctly before introducing them. An
				<span class="text-dark-emphasis">liquid ammonia kickstart - 2ppm</span> of ammonia added, we had readings !!
				Now the plan was to let the system be
				for a week, and query for bacteria ..`},
				{ImgPath: "", Txt: `Rather than waiting for tomato plants to germinate from seeds, we transplanted some from the farmhouse.
				These were <em>local varieties growing in a wasteland.</em> The primary objective was to observe the
				performance of the tomatoes and identify any operational challenges.`},
				{ImgPath: "", Txt: `March was huge success to me, we achieved the flow and siphon action as desired. Plus fish less cycling. We moved ahead with deep learning of the bell siphon and that I reckon would save a lot of time when expanding
				the setup further.`},
			},
		},

		{
			Header:   "April 2024",
			Slog:     "april-2024",
			CoverImg: "alienlettucefarm.png",
			Purport:  "Every step gives us much needed fillip in confidence. The system is not as responsive as would have loved it to be but we have some good water readings.",
			Gist:     `All that happened in April of 2024`,
			Author:   "Niranjan Awati",
			Owner:    "Eensymachines, Pune",
			Nav:      &BlogNav{Next: "may-2024", Prev: "march-2024", BaseHref: "/dear-diary"},
			Paras: []BlogPara{
				{ImgPath: "", Txt: `Ammonia had some initial inertia, as the readings would just drop with no signs of nitrites or nitrates. It  was a bit disheartening initially but we kept at it and gave it a slightly steeper kick start, Voila! a week down the line it was all working as expected. `},
				{ImgPath: "", Txt: `The tomato plants are thriving, exhibiting full, green foliage with no signs of vein or tip yellowing. Although we anticipated some bronzing of the leaves, this issue has not arisen. While some flowers have dropped prematurely, those that remain show no deficiencies.`},
				{ImgPath: "", Txt: `<table class="table table-hover">
				<thead>
					<tr>
						<th scope="col">Parameter</th>
						<th scope="col">Reading</th>
						<th scope="col">Remarks</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>Ammonia(NH3, NH4+)</td>
						<td>0.5ppm</td>
						<td>Ammonia was getting consumed as expected</td>
					</tr>
					<tr>
						<td>Nitrites (No2)</td>
						<td>3ppm</td>
						<td>yippe ! bacteria detected, water is potent</td>
					</tr>
					<tr>
						<td>Nitrate (No3)</td>
						<td>5ppm</td>
						<td>We are all set !- bacteria are at work and Nitrogen cycle is stable and established. </td>
					</tr>
					<tr>
						<td>pH</td>
						<td>7.6</td>
						<td>
							This isn't great news, but aquaponics is naturally acidifying process and we expect to see some  drop here later.  A higher pH while is ok for the fish, the nutrients that are available to plants at alkaline levels is recommended.
						</td>
					</tr>
				</tbody>
			</table>`},
				{ImgPath: "", Txt: `The bell siphon is functioning flawlessly without any leaks. The pump, however, is showing signs of wear. As a result, we will prioritize timing and automation, but only after replacing the current pump with a submersible one.`},
				{ImgPath: "", Txt: `Late in april we introduce Coi fish, black sharks - total of 15 fish, who found their new home a bit alien to start with but quickly adjusted to environment.`},
				{ImgPath: "", Txt: `Meanwhile Eensymachines has already started developing a prototype for the automation. `},
			},
		},
		{
			Header:   "May 2024",
			Slog:     "may-2024",
			CoverImg: "water_lettuce.png",
			Purport:  "When you have more doubts than you have answers, would it mean you are on the right path or atleast headed to one?",
			Gist:     `Developments that happened May of 2024`,
			Author:   "Niranjan Awati",
			Owner:    "Eensymachines, Pune",
			Nav:      &BlogNav{Next: "june-2024", Prev: "april-2024", BaseHref: "/dear-diary"},
			Paras: []BlogPara{
				{ImgPath: "", Txt: `May did not start on a great note - <br><br>
				It was a bit disheartening to see the aborting flowers, plus the early pre-monsoon showers caused more flowers to drop prematurely. Did we miss the finish line by a whisker ? `},
				{ImgPath: "", Txt: `Some internet tips proved helpful, particularly the advice that "tomatoes are self-pollinating flowers." We decided to stimulate pollination by gently poking the flowers with cotton buds, which was effective. Within a week, we had a dozen tomatoes hanging on the vines, with almost all the treated flowers bearing fruit. However, the challenges didn't end there. While we had mastered the nitrogen cycle, ensuring the fruits ripened without dropping prematurely required the addition of rock phosphate (phosphorus) and seaweed solution (potassium). I researched these additives online to confirm their safety for the fish.`},
				{ImgPath: "culture_lab_lettuce.png", Txt: `The system appeared to need an overhaul. Early rains had caused moss to proliferate in all corners of the tank, and I suspected significant growth on the exposed ends of the bell siphon. Although the plants showed no signs of distress or deficiencies, structural adjustments to the grow bed were necessary before the onset of the monsoons. The IMD had predicted an El Niño effect, with an expected 6% increase in monsoon rainfall for the region. I was uncertain if the tomato plants could withstand the gusty winds and heavy showers.`},
				{ImgPath: "", Txt: `What was significant then was we installed the automation electronics , with initial setting of 8 cycles / day, then bumped it up to 12 cycles per day operating at pulse every interval mode. Now we were not in danger of overworking the motor plus also had remote control on the flood - drain cycles.`},
				{ImgPath: "", Txt: `The siphon exhibited signs of air bleeding, with the snorkel showing an air leak, leading to slight water logging in the grow bed. My suspicion about moss obstructing the siphon flow was confirmed, or perhaps the sealant on the siphon joints had deteriorated.<br><br>All in all - May had raised more questions than it had answers. `},
				{ImgPath: "", Txt: `<table class="table table-hover">
				<thead>
					<tr>
						<th scope="col">Parameter</th>
						<th scope="col">Reading</th>
						<th scope="col">Remarks</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>Ammonia(NH3, NH4+)</td>
						<td>0ppm</td>
						<td>Some dead fish - the cycle is a bit imbalanced- do we risk adding ammonia when the fish in the system ?</td>
					</tr>
					<tr>
						<td>Nitrites (No2)</td>
						<td>2ppm</td>
						<td>Looks like we wait before we add the ammonia</td>
					</tr>
					<tr>
						<td>Nitrate (No3)</td>
						<td>5ppm</td>
						<td>No wonder the leaves look healthy, no distress</td>
					</tr>
					<tr>
						<td>pH</td>
						<td>7.6</td>
						<td>We had started to suspect the hardness of water the primary reason for this pH inertia.</td>
					</tr>
				</tbody>
			</table>`},
			},
		},
		{
			Header:   "June 2024",
			Slog:     "june-2024",
			CoverImg: "tomato_array.png",
			Purport:  "Knee jerk reactions are often riddled side effects, but unless you get to experience one you wouldn't learn",
			Gist:     `Unfoldingof events in June 2024`,
			Author:   "Niranjan Awati",
			Owner:    "Eensymachines, Pune",
			Nav:      &BlogNav{Next: "", Prev: "may-2024", BaseHref: "/dear-diary"},
			Paras: []BlogPara{
				{ImgPath: "liquid_chromatography.png", Txt: `Now that we had established the sustainability of the nitrogen cycle, we decided to undertake a system overhaul and cleaning. We had been experiencing consistent incidents of medium-sized Koi dying without any signs of infection or distress. Although aeration was adequate, we suspected the fish were succumbing to oxygen deprivation. Therefore, we decided to proceed with the system overhaul. It was a Herculean task to dispose of approximately 400 liters of water and move around 100 kg of gravel while ensuring no additional fish casualties.`},
				{ImgPath: "", Txt: `The tomato plants thrived post-overhaul, transitioning from green to red without any yellowing. However, the ammonia levels were not as encouraging, as the remaining fish, with minimal stocking, struggled to raise the ammonia levels sufficiently. To address this, we added more fish—ornamental sharks known for their hardiness and ability to tolerate water temperature fluctuations.`},
				{ImgPath: "", Txt: `Still have a feeling fish are falling short of generating that ammonia levels. The system is running low on a weak nitrogen cycle - but getting Tilapia is a challenge.  - The ammonia is just barely enough to sustain the plants we have used. We decide to let the cycle settle in - was the overhaul a little too early?`},
				{ImgPath: "", Txt: `<table class="table table-hover">
				<thead>
					<tr>
						<th scope="col">Parameter</th>
						<th scope="col">Reading</th>
						<th scope="col">Remarks</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>Ammonia(NH3, NH4+)</td>
						<td>0ppm</td>
						<td>System overhaul, no kickstart</td>
					</tr>
					<tr>
						<td>Nitrites (No2)</td>
						<td>0ppm</td>
						<td>Wasn't expecting any readings</td>
					</tr>
					<tr>
						<td>Nitrate (No3)</td>
						<td>0ppm</td>
						<td>Wasn't expecting any readings here too</td>
					</tr>
					<tr>
						<td>pH</td>
						<td>7.6</td>
						<td>duh !</td>
					</tr>
				</tbody>
			</table>`},
				{ImgPath: "", Txt: `Automation is doing just about fine, no issues there. We have a comfortable remote ssh connection to the device . It has saved us dozens of 34 km drives to & fro. A simple ssh from the comfort of the homes desk we can have a peep into the logs of the device - knowing the device will time exact 12 cycles of flood-drain in a day left us some bandwidth to think ahead of the impending problems`},
				{ImgPath: "", Txt: `June wasnt great per say, we are just hoping the system will pick up. Fish wouldnt die any longer, but we arent expecting the tomatoes to grow any faster than before.`},
			},
		},
	}
)
