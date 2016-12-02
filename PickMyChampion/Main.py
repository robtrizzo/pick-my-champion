from RiotAPI import RiotAPI

def main():
	api = RiotAPI("""<put your riot api key here>""")
	r = api.get_summoner_by_name('SpookyDaMoose')
	print(r)

if __name__ == "__main__":
	main()