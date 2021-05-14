import request from '@/utils/request'

export function groupList(gid) {
  return request({
    url: '/group/list',
    method: 'get',
    params: {gid: gid}
  })
}

export function groupAdd(data) {
  return request({
    url: '/group/add',
    method: 'post',
    data
  })
}

export function groupUpdate(data) {
  return request({
    url: '/group/update',
    method: 'post',
    data
  })
}

export function groupTestConnection(data) {
  console.log(data)
  return request({
    url: '/group/test_connection',
    method: 'post',
    data
  })
}

export function getConfigFromDB(data) {
  console.log(data)
  return request({
    url: '/group/get_config_from_db',
    method: 'post',
    data
  })
}

export function exportConfigToDB(data) {
  console.log(data)
  return request({
    url: '/group/export_config_to_db',
    method: 'post',
    data
  })
}

export function generateAppKeySecret() {
  return request({
    url: '/group/gen_app_key_secret',
    method: 'get',
  })
}

export function syncToProd(data) {
  console.log(data)
  return request({
    url: '/group/sync_to_prod',
    method: 'post',
    data
  })
}

export function exportRecord(data) {
  console.log(data)
  return request({
    url: '/group/export_record',
    method: 'post',
    data
  })
}

export function exportRecordContent(data) {
  console.log(data)
  return request({
    url: '/group/export_record_content',
    method: 'post',
    data
  })
}

export function exportRollback(data) {
  console.log(data)
  return request({
    url: '/group/export_rollback',
    method: 'post',
    data
  })
}