import { expect } from 'chai'
import { cloudRunService } from '../cloudRunService'
import { name } from '../globalConf'

describe('createCloudRunService', () => {
	let serviceName = name
	it('should have the correct name', (done: () => void) => {
		cloudRunService.name.apply(name => {
			expect(name).to.equal(serviceName)
			done()
		})
	})

	it('should have the correct ID', () => {
		expect(cloudRunService.id).to.exist
	})

	it('should have the correct URI', () => {
		expect(cloudRunService.uri).to.exist
	})
})
