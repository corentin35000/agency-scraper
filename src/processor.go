package main

import (
	"fmt"
	"time"
)

/**
 * RunScraper lance le scraping des annonces immobilières à intervalles réguliers.
 * Cette fonction est appelée depuis le point d'entrée de l'application.
 * @param {int} intervalMinutes - Intervalle de temps en minutes entre chaque cycle de scraping
 * return {void}
 */
func RunScraper(intervalMinutes int) {
	// Map globale pour suivre les références des biens déjà traités par les différentes agences
	processedReferencesAfedim := make(map[string]bool)
	processedReferencesGiboire := make(map[string]bool)
	processedReferencesFoncia := make(map[string]bool)
	processedReferencesAgenceDuColombier := make(map[string]bool)
	processedReferencesLaFrancaiseImmobiliere := make(map[string]bool)
	processedReferencesGuenno := make(map[string]bool)
	processedReferencesLaMotte := make(map[string]bool)
	processedReferencesKermarrec := make(map[string]bool)
	processedReferencesNestenn := make(map[string]bool)

	for {
		// Lancer le scraping pour l'agence Afedim
		processAgencyScraping(processedReferencesAfedim, "https://www.afedim.fr/fr/location/annonces/Appartement-Maison-Parking-Garage/Rennes-France/1-5-pieces/surface-0-100-m2/budget-0-90000-euros/rayon-10-km/disponible-/options-/exclusPlafondRess-/Resultats", "AFEDIM", "Afedim")

		// Lancer le scraping pour l'agence Giboire
		processAgencyScraping(processedReferencesGiboire, "https://www.giboire.com/recherche-location/appartement/?searchBy=default&address%5B%5D=RENNES&address%5B%5D=CHANTEPIE&address%5B%5D=CESSON+SEVIGNE&priceMax=700&nbBedrooms%5B%5D=1&transactionType%5B%5D=Location&searchBy=default", "GIBOIRE", "Giboire")

		// Lancer le scraping pour l'agence Foncia
		processAgencyScraping(processedReferencesFoncia, "https://fr.foncia.com/location/rennes-35--chantepie-35135--cesson-sevigne-35510/appartement?nbPiece=2--&prix=--700&advanced=", "FONCIA", "Foncia")

		// Lancer le scraping pour l'agence Agence du Colombier
		processAgencyScraping(processedReferencesAgenceDuColombier, "https://agenceducolombier.com/annonces/?filter_search_action%5B%5D=louer&filter_search_type%5B%5D=&nb-pieces=&min-chambres=&min-surface=&max-surface=&price_low=0&price_max=6000000&submit=LANCER+MA+RECHERCHE", "AGENCE DU COLOMBIER", "Agence du Colombier")

		// Lancer le scraping pour l'agence La Française Immobilière
		processAgencyScraping(processedReferencesLaFrancaiseImmobiliere, "https://www.la-francaise-immobiliere.fr/location/?post_types=location&categorie%5B%5D=27&zone%5B%5D=6212&zone%5B%5D=6204&zone%5B%5D=6214&nb_chambres_min=0&nb_chambres_max=&prix_min=0&prix_max=700&submitted=1&o=date-desc&action=load_search_results&wia_6_type=&searchOnMap=0&wia_1_reference=", "LA FRANCAISE IMMOBILIERE", "La Française Immobilière")

		// Lancez le scraping pour l'agence Guenno
		processAgencyScraping(processedReferencesGuenno, "https://www.guenno.com/biens/recherche?mandate_type=2&realty_type%5B%5D=1&number_room%5B%5D=2&min_surface=&town=RENNES+35000&price_max=700", "GUENNO", "Guenno")

		// Lancer le scraping pour l'agence La Motte
		processAgencyScraping(processedReferencesLaMotte, "https://www.kermarrec-habitation.fr/location/?post_type=location&false-select=on&99795fbc=&ville%5B%5D=cesson-sevigne-35510&ville%5B%5D=chantepie-35135&ville%5B%5D=rennes-35000&typebien%5B%5D=appartement&budget_max=700&reference=&rayon=0&avec_carte=false&tri=pertinence", "LA MOTTE", "La Motte")

		// Lancer le scraping pour l'agence Kermarrec
		processAgencyScraping(processedReferencesKermarrec, "https://www.kermarrec-habitation.fr/location/?post_type=location&false-select=on&99795fbc=&ville%5B%5D=cesson-sevigne-35510&ville%5B%5D=chantepie-35135&ville%5B%5D=rennes-35000&typebien%5B%5D=appartement&budget_max=700&reference=&rayon=0&avec_carte=false&tri=pertinence", "KERMARREC", "Kermarrec")

		// Lancer le scraping pour l'agence Nestenn
		processAgencyScraping(processedReferencesNestenn, "https://immobilier-rennes-centre.nestenn.com/?action=listing&prestige=0&meuble=0&transaction=louer&list_ville=35+Rennes%2C35135+Chantepie%2C35510+Cesson-S%C3%A9vign%C3%A9&list_type=Appartement&type=Appartement&prix_max=700&pieces=2", "NESTENN", "Nestenn")

		// Attendre avant le prochain cycle
		time.Sleep(time.Duration(intervalMinutes) * time.Minute)
	}
}

/**
 * processAgencyScraping lance le scraping pour une agence immobilière spécifique.
 * @param {map[string]bool} processedReferences - Map contenant les références des biens déjà traités.
 * @param {string} url - L'URL de la page de l'agence à scraper.
 * @param {string} titleMessageTelegram - Le titre du message Telegram.
 * @param {Agency} nameAgency - Le nom de l'agence.
 * @return {void}
 */
func processAgencyScraping(processedReferences map[string]bool, url string, titleMessageTelegram string, nameAgency Agency) {
	// Créer une nouvelle instance de CollyService
	collyService := NewCollyService()

	// Récupérer les annonces complètes depuis l'agence
	newAnnouncements := collyService.ScrapeAnnouncement(nameAgency, url)

	// Comparer les références des biens pour détecter les nouvelles annonces
	for _, announcement := range newAnnouncements {
		if _, exists := processedReferences[announcement.propertyReference]; !exists {
			// Nouvelle annonce détectée
			fmt.Println("Nouvelle annonce détectée référence :", announcement.propertyReference)
			processedReferences[announcement.propertyReference] = true

			// Envoie un message sur le canal Telegram
			sendTelegramMessageToPublicChannel(fmt.Sprintf(
				"%s\nNouvelle annonce immobilière !\nRéférence : %s\nURL : %s",
				titleMessageTelegram,
				announcement.propertyReference,
				announcement.url,
			))
		}
	}
}
