import * as pulumi from '@pulumi/pulumi'
import * as gcp from '@pulumi/gcp'
import { project, region, event, labels } from './globalConf'
import { serviceAccount } from './serviceAccount'
import { cloudRunService } from './cloudRunService'

interface Events {
	name: string
	location: string
	id: pulumi.Output<string>
}

const events: Events[] = []

event.buckets.forEach(bucket => {
	const trigger = new gcp.eventarc.Trigger(bucket.name, {
		name: `${bucket.name}-scan`,
		project: project,
		location: bucket.location,
		destination: {
			cloudRunService: {
				service: cloudRunService.name,
				path: '/scan',
				region: region
			}
		},
		matchingCriterias: [
			{
				attribute: 'type',
				value: event.type
			},
			{
				attribute: 'bucket',
				value: bucket.name
			}
		],
		serviceAccount: serviceAccount.email,
		labels: labels
	})

	const scanBucketIAM = new gcp.storage.BucketIAMMember(bucket.name, {
		bucket: bucket.name,
		role: 'roles/storage.objectUser',
		member: serviceAccount.member
	})
	events.push({ name: `${bucket.name}-scan`, location: bucket.location, id: trigger.id })
})

export const eventTrigger = events
