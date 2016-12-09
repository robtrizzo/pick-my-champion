
class ChampDatabase(object):

	def __init__(self, champInfo):
		self.champInfo = champInfo

	def populateChampMatrix(self, champInfo):
		champMatrix = [[0 for x in range(len(champInfo['data']))] for y in range(2)]
		i = 0
		name = ''
		id = 0
		for key in champInfo['data']:
			name = key
			id = champInfo['data'][name]['id'] 
			champMatrix[0][i] = id
			champMatrix[1][i] = name
			print('ID: %d   Name: %s' % (champMatrix[0][i], champMatrix[1][i]))
			i+=1
		return champMatrix
