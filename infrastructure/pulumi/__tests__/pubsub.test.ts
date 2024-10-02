import { expect } from 'chai'
import { pubSub } from '../pubsub'

describe('createPubSub', () => {
	it('should have service account created', () => {
		expect(pubSub.topic).to.exist
		expect(pubSub.subscription).to.exist
		expect(pubSub.deadLetterTopic).to.exist
	})
})
