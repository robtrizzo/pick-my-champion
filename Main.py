from RiotAPI import RiotAPI

def main():
	api = RiotAPI('RGAPI-a385434d-2a02-40fe-8fdb-df67d38a506f')
	r = api.get_summoner_by_name('SpookyDaMoose')
	print(r)

if __name__ == "__main__":
	main()