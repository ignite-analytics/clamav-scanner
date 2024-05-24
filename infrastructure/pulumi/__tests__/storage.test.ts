import { expect } from 'chai'
import { storageBucket } from '../storage'

describe('createStorage', () => {
	it('should have storage buckets created', () => {
		expect(storageBucket.length).to.equal(2)
	})
})
