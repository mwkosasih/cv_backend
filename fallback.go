package main

import "strings"

// GetFallbackReply returns a rich profile response when Gemini is offline or unconfigured
func GetFallbackReply(msg string) string {
	lower := strings.ToLower(msg)

	if reply, ok := matchStockbit(lower); ok {
		return reply
	}
	if reply, ok := matchSkills(lower); ok {
		return reply
	}
	if reply, ok := matchCertifications(lower); ok {
		return reply
	}
	if reply, ok := matchContact(lower); ok {
		return reply
	}
	if reply, ok := matchQasir(lower); ok {
		return reply
	}
	if reply, ok := matchEducation(lower); ok {
		return reply
	}
	if reply, ok := matchCareer(lower); ok {
		return reply
	}
	if reply, ok := matchAnime(lower); ok {
		return reply
	}
	if reply, ok := matchGaming(lower); ok {
		return reply
	}
	if reply, ok := matchBirthplace(lower); ok {
		return reply
	}
	if reply, ok := matchLifeDream(lower); ok {
		return reply
	}
	if reply, ok := matchHobbies(lower); ok {
		return reply
	}

	return defaultIntro()
}

func matchStockbit(lower string) (string, bool) {
	if strings.Contains(lower, "stockbit") || strings.Contains(lower, "idx") || strings.Contains(lower, "exchange") || strings.Contains(lower, "securities") || strings.Contains(lower, "intraservices") || strings.Contains(lower, "balance") {
		return "At **Stockbit** (May 2021 – Present, 5+ years • Indonesia), I work as a Back End Developer handling critical financial infrastructure.\n\n" +
			"Notable Activities & Contributions:\n" +
			"• **PHP to Go Migration**: Successfully migrated legacy backend services from PHP to Go, improving system maintainability, operational efficiency, and overall throughput.\n" +
			"• **Monolith to Microservices Transition**: Contributed strategically to decomposing monolithic services into microservices, enabling better service isolation, higher availability, and streamlined team workflows.\n" +
			"• **Direct IDX Price Feed Integration**: Migrated the real-time stock price data feed from third-party vendor platforms to a direct, low-latency integration with the Indonesia Stock Exchange (IDX).\n" +
			"• **Stockbit Intraservices System Migration**: Migrated the Stockbit Securities Intraservices system from an external third-party vendor solution to an internally developed, maintained, and optimized system.\n" +
			"• **User Balance Management**: Developed and maintained critical backend services responsible for user balance calculations, deposit processing, withdrawal processing, and balance adjustments with strict transactional consistency.\n\n" +
			"Technologies Used: Golang, Kubernetes, AWS, PostgreSQL, TimescaleDB, Redis, Kafka, NATS, SQS, WebSockets, Git, Grafana.", true
	}
	return "", false
}

func matchSkills(lower string) (string, bool) {
	if strings.Contains(lower, "skill") || strings.Contains(lower, "tech") || strings.Contains(lower, "stack") || strings.Contains(lower, "tool") {
		return "My technical expertise centers around **Golang** and scalable backend architectures:\n\n" +
			"• **Languages & Core**: Golang (Go), PHP, JavaScript, SQL\n" +
			"• **Architecture & Design**: Microservices, System Design, Database Design, REST & gRPC APIs, Distributed Systems\n" +
			"• **Databases & Caching**: PostgreSQL, TimescaleDB, Redis, MySQL\n" +
			"• **Messaging & Streaming**: Apache Kafka, NATS, AWS SQS, WebSockets, GCP Pub/Sub\n" +
			"• **Cloud & DevOps**: Kubernetes, AWS, GCP, Git\n" +
			"• **Observability & AI**: Grafana, NewRelic, LLM Integrations", true
	}
	return "", false
}

func matchCertifications(lower string) (string, bool) {
	if strings.Contains(lower, "cert") || strings.Contains(lower, "hackerrank") || strings.Contains(lower, "sertifikat") {
		return "Here are my verified **HackerRank Certifications**:\n\n" +
			"• **Rest API (Intermediate) Certificate**\n" +
			"  https://www.hackerrank.com/certificates/87c5a66ec989\n\n" +
			"• **Problem Solving (Basic) Certificate**\n" +
			"  https://www.hackerrank.com/certificates/9e06391dfb4e", true
	}
	return "", false
}

func matchContact(lower string) (string, bool) {
	if strings.Contains(lower, "contact") || strings.Contains(lower, "email") || strings.Contains(lower, "phone") || strings.Contains(lower, "hire") || strings.Contains(lower, "reach") || strings.Contains(lower, "linkedin") {
		return "I'd love to connect! You can reach me directly via:\n\n" +
			"• **Email**: mwkosasih@gmail.com\n" +
			"• **Mobile / WhatsApp**: +62 856-9366-9469\n" +
			"• **LinkedIn**: [linkedin.com/in/mohammad-wildan-kosasih](https://www.linkedin.com/in/mohammad-wildan-kosasih)\n" +
			"• **Location**: Indonesia (Open to remote and hybrid opportunities)", true
	}
	return "", false
}

func matchQasir(lower string) (string, bool) {
	if strings.Contains(lower, "qasir") || strings.Contains(lower, "pos") {
		return "At **Qasir.id** (June 2019 – May 2021, 2 years • Indonesia), I served as a Back End Developer:\n\n" +
			"• Built and maintained core transaction and inventory features for the Qasir POS mobile and web platforms.\n" +
			"• Revamped the monolithic API architecture into scalable microservices and led the transition of legacy PHP endpoints to Go.\n\n" +
			"Technologies Used: Golang, PHP, GCP, PostgreSQL, MySQL, Redis, Kafka, Cloud Pub/Sub, Git, NewRelic.", true
	}
	return "", false
}

func matchEducation(lower string) (string, bool) {
	if strings.Contains(lower, "education") || strings.Contains(lower, "university") || strings.Contains(lower, "degree") || strings.Contains(lower, "pasim") {
		return "I hold an **Associate's Degree (A.Md.), Management Information Systems, General** from **PASIM National University** (2013 – 2016).\n\n" +
			"I also hold certifications in **Rest API (Intermediate)** and **Problem Solving (Basic)**.", true
	}
	return "", false
}

func matchCareer(lower string) (string, bool) {
	if strings.Contains(lower, "rajamobil") || strings.Contains(lower, "edi") || strings.Contains(lower, "surya") || strings.Contains(lower, "experience") || strings.Contains(lower, "career") {
		return "Here is a quick overview of my engineering career:\n\n" +
			"1. **Stockbit** (May 2021 – Present): Back End Developer focusing on balance ledger, deposit/withdrawal processing, IDX direct feed, and intraservices migration.\n" +
			"2. **Qasir.id** (June 2019 – May 2021): Back End Developer revamping POS monolith to Go microservices.\n" +
			"3. **RajaMobil.com** (Feb 2018 – June 2019): Full Stack Developer building MY.Rajamobil.com and mobibagus.com management panels.\n" +
			"4. **PT EDI Indonesia** (Jan 2017 – Feb 2018): Full Stack Developer building enterprise Procurement Request systems for the company.\n" +
			"5. **Surya Putra Sukses** (Sep 2014 – Jan 2017): Full Stack Developer developing clinical, retail, school, HRMS, procurement, and tax payment systems.", true
	}
	return "", false
}

func matchAnime(lower string) (string, bool) {
	if strings.Contains(lower, "anime") || strings.Contains(lower, "manga") || strings.Contains(lower, "nonton") || strings.Contains(lower, "watch") {
		return "Outside of engineering, I'm a huge **anime fan**! 🎬\n\n" +
			"Some of my all-time favorites include:\n" +
			"• **Shonen & Action**: *One Piece, Naruto, Bleach, Hunter x Hunter, Attack on Titan, Jujutsu Kaisen, Demon Slayer, Solo Leveling, Vinland Saga*\n" +
			"• **Fantasy**: *Frieren: Beyond Journey's End, Mushoku Tensei*\n" +
			"• **Drama & Mystery**: *Death Note, Erased, A Silent Voice (Koe no Katachi)*\n" +
			"• **Comedy & Slice of Life**: *Kaguya-sama: Love is War, KonoSuba, Gintama, The Pet Girl of Sakurasou, Nisekoi*\n\n" +
			"Always open to great anime recommendations!", true
	}
	return "", false
}

func matchGaming(lower string) (string, bool) {
	if strings.Contains(lower, "game") || strings.Contains(lower, "gaming") || strings.Contains(lower, "dota") || strings.Contains(lower, "terraria") || strings.Contains(lower, "stardew") || strings.Contains(lower, "main") {
		return "When I'm not writing code, I love gaming to unwind and strategize! 🎮\n\n" +
			"• **Dota 2**: My most-played game by far — love the competitive team dynamics and strategic depth.\n" +
			"• **Terraria**: Love the open sandbox exploration, boss battles, and crafting progression.\n" +
			"• **Stardew Valley**: Perfect cozy game for relaxing and farming.\n\n" +
			"Got a favorite game or down for a Dota match sometime?", true
	}
	return "", false
}

func matchBirthplace(lower string) (string, bool) {
	if strings.Contains(lower, "cirebon") || strings.Contains(lower, "born") || strings.Contains(lower, "birth") || strings.Contains(lower, "lahir") || strings.Contains(lower, "1995") || strings.Contains(lower, "asal") {
		return "I was born in **1995** in **Cirebon Regency (Kabupaten Cirebon), West Java, Indonesia**, and I'm currently based in Indonesia.", true
	}
	return "", false
}

func matchLifeDream(lower string) (string, bool) {
	if strings.Contains(lower, "dream") || strings.Contains(lower, "travel") || strings.Contains(lower, "keliling dunia") || strings.Contains(lower, "family") || strings.Contains(lower, "keluarga") || strings.Contains(lower, "cita") {
		return "One of my biggest personal dreams and life aspirations is to **travel around the world together with my family** ✈️🌍, discovering new cultures, scenic landscapes, and global cuisines along the way!", true
	}
	return "", false
}

func matchHobbies(lower string) (string, bool) {
	if strings.Contains(lower, "hobby") || strings.Contains(lower, "hobbies") || strings.Contains(lower, "hobi") || strings.Contains(lower, "personal") || strings.Contains(lower, "free time") || strings.Contains(lower, "downtime") {
		return "Outside of software engineering, here's a bit about me:\n\n" +
			"• **Anime**: Big anime fan (watching everything from *Frieren*, *One Piece*, and *Bleach* to *Kaguya-sama* and *Vinland Saga*).\n" +
			"• **Gaming**: Play a lot of **Dota 2**, plus sandbox favorites like **Terraria** and **Stardew Valley**.\n" +
			"• **Life Goal**: Work hard toward my dream of traveling the world with my family.", true
	}
	return "", false
}

func defaultIntro() string {
	return "Hello! I am **Mohammad Wildan Kosasih's AI Twin**.\n\n" +
		"I am a Backend Engineer specializing in Go (Golang), building scalable, reliable, and high-concurrency systems.\n\n" +
		"You can ask me about my engineering work at **Stockbit** or **Qasir.id**, my **system design** projects, or even my hobbies like **anime**, **gaming (Dota 2)**, and life outside of tech!"
}
