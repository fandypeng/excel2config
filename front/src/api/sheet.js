import request from '@/utils/request'

export function getList(params) {
  return request({
    url: '/excel/list',
    method: 'get',
    params
  })
}

export function deleteExcel(params) {
  return request({
    url: '/excel/delete',
    method: 'post',
    params
  })
}


export function addExcel(params) {
  return request({
    url: '/excel/create',
    method: 'post',
    params
  })
}


export function updateExcel(params) {
  return request({
    url: '/excel/update',
    method: 'post',
    params
  })
}

export function getSheetList(params) {
  return request({
    url: '/excel/sheet_list',
    method: 'post',
    params
  })
}



export function exportExcel(params) {
  return request({
    url: '/excel/export',
    method: 'post',
    params
  })
}

export function exportExcelProd(params) {
  return request({
    url: '/excel/export_prod',
    method: 'post',
    params
  })
}






