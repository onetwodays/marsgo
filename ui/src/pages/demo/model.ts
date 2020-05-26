import {Effect, Reducer,Subscription,request} from 'umi'

export interface DemoTable{
    id:number;
    name:string;
    desc:string;
    url:string;
}

export interface DemoModelStateType{
    rows:DemoTable[];
    filterKey:string;
}


export interface DemoModelType{
    namespace:"demo";
    state:DemoModelStateType;
    reducers:{
        save:Reducer<DemoModelStateType>,
    };
    effects:{
        fetch:Effect,
    };
    subscriptions:{
        setup:Subscription,
    };
}

const DemoModel: DemoModelType= {
    namespace:'demo',
    state:{
        rows:[],
        filterKey:""
    },
    reducers:{
        save(state,action){return { ...state,...action.payload}}
    },


    effects:{
        *fetch({type,payload},{put,call,select}){
            yield put({type:'save',payload:payload});

        }

    },
    subscriptions:{
        setup({dispatch,history}){
            return history.listen(({pathname,query})=>{
                if(pathname==='/demo'){
                    dispatch({
                        type:'fetch',
                        payload:{
                            rows:[{id:1,name:'caoxue',desc:'曹雪',url:''}],
                            filterKey:"caoxue"
                        }
                    });
                }
            });
        }
    },
};
