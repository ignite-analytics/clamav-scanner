import { expect } from 'chai'
import { eventTrigger } from '../eventarch'

describe('createEventTrigger', () => {
	it('should have event triegger created', () => {
		expect(eventTrigger.length).to.be.greaterThan(0)
	})
})
