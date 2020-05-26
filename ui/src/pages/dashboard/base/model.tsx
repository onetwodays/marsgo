import {Effect , Reducer,Subscription,request} from 'umi'

export interface BaseModelStateType{
    name: string;
}

export interface BaseModelType{
    namespace:'base';
    state:BaseModelStateType;
    effects:{
       query:Effect; 
    };
    reducers:{
        save: Reducer< BaseModelStateType >;
    };
    subscriptions:{
        setup:Subscription
    };

}

const BaseModel : BaseModelType = {
    namespace:'base',

    state:{
        name:'hello-base',
    },

    effects:{
        *query({payload},{call,put}){

        },
        *fetch({type,payload},{call,put,select}){
            yield put({
                type:'save',
                payload:{
                    name:payload.name||'zhouhao',
                }

            });
        },

        *fetchHero({type,payload},{call,put,select}){
            const data = yield request('/apimock/web201605/js/herolist.json');
            const localData = [
                {
                ename: 105,
                cname: '廉颇',
                title: '正义爆轰',
                new_type: 0,
                hero_type: 3,
                skin_name: '正义爆轰|地狱岩魂',
                },
                {
                ename: 106,
                cname: '小乔',
                title: '恋之微风',
                new_type: 0,
                hero_type: 2,
                skin_name: '恋之微风|万圣前夜|天鹅之梦|纯白花嫁|缤纷独角兽',
                },
            ];

            yield put({
                type:'save',
                payload:{
                    name:JSON.stringify(data  || localData),
                }

            });


        },
    },

    reducers:{
        save(state,action){
            return {
                ...state,
                ...action.payload,
            };
        },
    },

    subscriptions:{
        setup({dispatch,history}){
            return history.listen(({pathname,query})=>{
                if(pathname==='/dashboard/base'){
                    dispatch({
                        type:'fetchHero'
                    })
                }

            });

        },
    }



};
export default BaseModel;