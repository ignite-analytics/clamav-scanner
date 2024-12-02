import * as pulumi from '@pulumi/pulumi'
import * as gcp from '@pulumi/gcp'
import { name, schedule } from './globalConf'
import { serviceAccount } from './serviceAccount'
import { cloudRunService } from './cloudRunService'

const createCloudRunScheduler = new gcp.cloudscheduler.Job('default', {
	name: `${name}-update`,
	description: 'ClamAV Scanner scheduled update job',
	region: schedule.location,
	schedule: schedule.cron,
	timeZone: schedule.timezone,
	attemptDeadline: '320s',
	retryConfig: {
		retryCount: 1
	},
	httpTarget: {
		httpMethod: 'POST',
		uri: pulumi.interpolate`${cloudRunService.uri}/update`,
		oidcToken: {
			serviceAccountEmail: serviceAccount.email,
			audience: cloudRunService.uri
		}
	}
})

export const cloudScheduler = {
	id: createCloudRunScheduler.id,
	name: createCloudRunScheduler.name
}
