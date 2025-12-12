import * as gcp from '@pulumi/gcp'
import { project, pubsub, event, labels } from './globalConf'
import { serviceAccount } from './serviceAccount'

const topic = new gcp.pubsub.Topic(`${pubsub.topic}`, {
	name: pubsub.topic,
	labels: labels
})

const deadLetterTopic = new gcp.pubsub.Topic(`${pubsub.topic}-dead-letter`, {
	name: `${pubsub.topic}-dead-letter`,
	labels: labels
})

const deadLetterTopicIam = new gcp.pubsub.TopicIAMMember(`${deadLetterTopic.name}`, {
	topic: deadLetterTopic.name,
	role: 'roles/pubsub.publisher',
	member: `serviceAccount:${pubsub.gcsServiceAccount}`
})

const subscription = new gcp.pubsub.Subscription(
	`${pubsub.subscription}`,
	{
		name: pubsub.subscription,
		topic: topic.name,
		ackDeadlineSeconds: 20,
		retainAckedMessages: true,
		messageRetentionDuration: '1200s',
		deadLetterPolicy: {
			deadLetterTopic: deadLetterTopic.id,
			maxDeliveryAttempts: 5
		},
		labels: labels
	},
	{
		dependsOn: [topic, deadLetterTopic]
	}
)

const subscriptionDeadLetter = new gcp.pubsub.Subscription(
	`${pubsub.subscription}-dead-letter`,
	{
		name: `${pubsub.subscription}-dead-letter`,
		topic: deadLetterTopic.name,
		ackDeadlineSeconds: 20,
		retryPolicy: {
			minimumBackoff: '1s',
			maximumBackoff: '10s'
		},
		labels: labels
	},
	{
		dependsOn: [deadLetterTopic]
	}
)

const topics = [topic.name, deadLetterTopic.name]

topics.forEach(topicName => {
	topicName.apply(name => {
		new gcp.pubsub.TopicIAMMember(name, {
			project: project,
			topic: name,
			role: 'roles/pubsub.editor',
			member: serviceAccount.member
		})
	})
})

const subscriptionIamBinding = new gcp.pubsub.SubscriptionIAMBinding('subscribers-binding', {
	subscription: subscription.name,
	role: 'roles/pubsub.subscriber',
	members: [...pubsub.subscribers, pubsub.gcsServiceAccount].map(subscriber => `serviceAccount:${subscriber}`)
})

export const pubSub = {
	topic: topic.name,
	subscription: subscription.name,
	deadLetterTopic: deadLetterTopic.name
}
