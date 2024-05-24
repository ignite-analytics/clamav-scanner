import { expect } from 'chai'
import { cloudScheduler } from '../cloudScheduler'
import { name } from '../globalConf'

describe('createCloudRunScheduler', () => {
	let scheduleName = `${name}-update`

	it('should have the correct name', done => {
		cloudScheduler.name.apply(name => {
			try {
				expect(name).to.equal(scheduleName)
				done()
			} catch (error) {
				done(error)
			}
		})
	})

	it('should have the correct ID', () => {
		expect(cloudScheduler.id).to.exist
	})
})
