package etcd

func (core *CoreEntity) WithLeaseRetryBefore(handle func()) {
	core.leaseRetryBefore = handle
}
func (core *CoreEntity) WithLeaseRetryAfter(handle func()) {
	core.leaseRetryAfter = handle
}
