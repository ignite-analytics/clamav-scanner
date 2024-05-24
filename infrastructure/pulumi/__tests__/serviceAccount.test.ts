import { expect } from 'chai'
import { serviceAccount } from '../serviceAccount'

describe('createServiceAccount', () => {
	it('should have service account created', () => {
		expect(serviceAccount.email).to.exist
		expect(serviceAccount.id).to.exist
		expect(serviceAccount.member).to.exist
	})
})
