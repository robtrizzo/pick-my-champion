from RiotAPI import RiotAPI
from ChampDatabase import ChampDatabase

def main():
	api = RiotAPI("""< Riot API Key >""")
	sn = api.get_summoner_by_name('SpookyDaMoose')
	print(sn)
	r = api.get_champion_ids()
	database = ChampDatabase(r)
	matrix = database.populateChampMatrix(r)

if __name__ == "__main__":
	main()