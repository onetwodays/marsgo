import request from '../../../util/request';

export function queryList() {
  return request('/apimock/cards');
}

export function deleteOne(id) {
  return request(`/apimock/cards/${id}`, {
    method: 'DELETE'
  });
}

export function addOne(data) {
  return request('/apimock/cards/add', {
    headers: {
      'content-type': 'application/json',
    },
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function getStatistic(id) {
  return request(`/apimock/cards/${id}/statistic`);
}
