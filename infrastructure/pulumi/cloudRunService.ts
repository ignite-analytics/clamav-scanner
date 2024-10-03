import * as gcp from '@pulumi/gcp'
import { name, project, region, image, service, storage, labels, pubsub } from './globalConf'
import { serviceAccount } from './serviceAccount'

const createCloudRunService = new gcp.cloudrunv2.Service('default', {
	name: name,
	project: project,
	location: region,
	ingress: 'INGRESS_TRAFFIC_INTERNAL_ONLY',
	template: {
		serviceAccount: serviceAccount.email,
		executionEnvironment: 'EXECUTION_ENVIRONMENT_GEN2',
		containers: [
			{
				image: `${image.repository}:${image.tag}`,
				resources: {
					startupCpuBoost: true,
					limits: {
						cpu: '1',
						memory: '4Gi'
					}
				},
				ports: {
					containerPort: service.port
				},
				envs: [
					{
						name: 'MIRROR_BUCKET',
						value: storage.mirrorBucket
					},
					{
						name: 'QUARANTINE_BUCKET',
						value: storage.quarantineBucket
					},
					{
						name: 'LISTEN_ADDRESS',
						value: ':1337'
					},
					{
						name: 'PROJECT_ID',
						value: project
					},
					{
						name: 'PUBSUB_TOPIC',
						value: pubsub.topic
					}
				]
			}
		],
		maxInstanceRequestConcurrency: 20,
		scaling: {
			minInstanceCount: 1,
			maxInstanceCount: 5
		}
	},
	traffics: [
		{
			type: 'TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST',
			percent: 100
		}
	],
	labels: labels
})

const serviceAccountIAM = new gcp.cloudrunv2.ServiceIamMember('default', {
	name: createCloudRunService.name,
	project: project,
	location: region,
	role: 'roles/run.invoker',
	member: serviceAccount.member
})

export const cloudRunService = {
	name: createCloudRunService.name,
	id: createCloudRunService.id,
	uri: createCloudRunService.uri
}
