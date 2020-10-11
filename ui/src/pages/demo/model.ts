import { Effect, Reducer, Subscription, request } from 'umi'


export interface rowProps {
    id: number;
    name: string;
    desc: string;
    url: string;
}

export interface DemoModelState {
    rows: rowProps[];
    filterKey: string;
}


export interface DemoModelType {
    namespace: "demo";
    state: DemoModelState;
    reducers: {
        save: Reducer<DemoModelState>;
    };
    effects: {
        fetch: Effect;
    };
    subscriptions: {
        setup: Subscription;
    };
}

const DemoModel: DemoModelType = {
    namespace: 'demo',
    state: {
        rows: [],
        filterKey: "",
    },
    reducers: {
        save(state, { payload }) { return { ...state, ...payload }; }
    },


    effects: {
        *fetch({ type, payload }, { put, call, select }) {
            yield put({ type: 'save', payload: payload });

        }

    },
    subscriptions: {
        setup({ dispatch, history }) {
            return history.listen(({ pathname, query }) => {
                if (pathname === '/demo') {
                    dispatch({
                        type: 'fetch',
                        payload: {
                            rows: [{ id: 1, name: 'caoxue', desc: '曹雪', url: '' }],
                            filterKey: "caoxue"
                        }
                    });
                }
            });
        }
    },
};

export default DemoModel;
