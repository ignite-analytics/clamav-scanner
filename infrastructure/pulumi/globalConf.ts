export const name = 'clamav-scanner' // service name
export const project = '' // GCP project ID
export const region = '' // service location
export const image = {
	repository: 'europe-docker.pkg.dev/clamav-scanner/clamav-scanner',
	tag: process.env['VERSION'] || 'latest'
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
