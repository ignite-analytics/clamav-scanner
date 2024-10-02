import * as pulumi from '@pulumi/pulumi'
import * as gcp from '@pulumi/gcp'
import { name, project, storage, labels, kms } from './globalConf'
import { serviceAccount } from './serviceAccount'

export const gcsAccount = gcp.storage.getProjectServiceAccount({
	project: project
})

const clamavKeyring = new gcp.kms.KeyRing('default', {
	name: `${name}-keyring`,
	location: kms.location,
	project: project
})

const clamavCryptoKey = new gcp.kms.CryptoKey(
	'default',
	{
		name: `${name}-key`,
		keyRing: clamavKeyring.id,
		rotationPeriod: kms.rotationPeriod
	},
	{
		dependsOn: [clamavKeyring]
	}
)

const gcsAccountKmsIamMember = new gcp.kms.CryptoKeyIAMMember(
	'default',
	{
		cryptoKeyId: clamavCryptoKey.id,
		role: 'roles/cloudkms.cryptoKeyEncrypterDecrypter',
		member: gcsAccount.then(gcsAccount => `serviceAccount:${gcsAccount.emailAddress}`)
	},
	{
		dependsOn: [clamavCryptoKey]
	}
)

const gcsAccountPubsubIamMember = new gcp.projects.IAMMember('pubsub', {
	project: project,
	role: 'roles/pubsub.publisher',
	member: gcsAccount.then(gcsAccount => `serviceAccount:${gcsAccount.emailAddress}`)
})

const createStorage = [storage.mirrorBucket, storage.quarantineBucket]

const bucketIds: pulumi.Output<string>[] = []

createStorage.forEach(bucket => {
	const clamavBucket = new gcp.storage.Bucket(
		bucket,
		{
			name: bucket,
			location: 'EU',
			publicAccessPrevention: 'enforced',
			uniformBucketLevelAccess: true,
			versioning: {
				enabled: true
			},
			logging: {
				logBucket: storage.logBucket
			},
			encryption: {
				defaultKmsKeyName: clamavCryptoKey.id
			},
			labels: labels
		},
		{
			dependsOn: [clamavCryptoKey]
		}
	)

	const clamavBucketIAM = new gcp.storage.BucketIAMMember(
		bucket,
		{
			bucket: clamavBucket.name,
			role: 'roles/storage.objectUser',
			member: serviceAccount.member
		},
		{
			dependsOn: [clamavBucket]
		}
	)
	bucketIds.push(clamavBucket.id)
})

export const storageBucket = bucketIds
