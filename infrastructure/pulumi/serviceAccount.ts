import * as gcp from '@pulumi/gcp'
import { name, project } from './globalConf'

const createServiceAccount = new gcp.serviceaccount.Account('default', {
	accountId: name,
	project: project,
	displayName: 'ClamAV Scanner',
	description: 'ClamAV Scanner Service Account',
	createIgnoreAlreadyExists: true,
	disabled: false
})

const serviceAccountIAM = new gcp.projects.IAMMember(
	'default',
	{
		project: project,
		role: 'roles/eventarc.eventReceiver',
		member: createServiceAccount.member
	},
	{
		dependsOn: [createServiceAccount]
	}
)

export const serviceAccount = {
	member: createServiceAccount.member,
	email: createServiceAccount.email,
	id: createServiceAccount.id
}
