export const name = 'clamav-scanner' // service name
export const project = '' // GCP project ID
export const region = '' // service location
export const image = {
	repository: 'ghcr.io/ignite-analytics/clamav-scanner',
	tag: '0.1.0'
}
export const schedule = {
	cron: '37 */2 * * *',
	timezone: 'Europe/Oslo',
	location: '' // cheduler job location
}
export const service = {
	port: 1337
}
export const event = {
	type: 'google.cloud.storage.object.v1.finalized',
	buckets: [] // [{name: 'bucket-name', location: 'EU'}]
}
export const storage = {
	location: 'EU',
	mirrorBucket: `${name}-mirror`,
	quarantineBucket: `${name}-quarantine`,
	logBucket: '' // where to store bucket logs
}
export const kms = {
	location: 'europe', // KMS keyring location
	rotationPeriod: '7776000s' // 90 days
}
export const labels = {
	manager: 'pulumi',
	service: name
}
