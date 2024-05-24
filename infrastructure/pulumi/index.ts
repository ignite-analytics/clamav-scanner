import { serviceAccount } from './serviceAccount'
import { storageBucket } from './storage'
import { cloudRunService } from './cloudRunService'
import { cloudScheduler } from './cloudScheduler'
import { eventTrigger } from './eventarch'

export const clamav = {
	serviceAccount: serviceAccount,
	service: cloudRunService,
	schedule: cloudScheduler,
	events: eventTrigger,
	storage: storageBucket
}
