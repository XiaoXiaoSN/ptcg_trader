import http from "k6/http";
import { check } from "k6";

let randItemID = () => {
    const itemList = [1, 2, 3, 4]
    return itemList[Math.floor(Math.random() * itemList.length)]
}

let randOrderType = () => {
    const orderTypeList = [1, 2]
    Math.floor(Math.random() * orderTypeList.length)
    return orderTypeList[Math.floor(Math.random() * orderTypeList.length)]
}

let randPrice = () => {
    return (Math.random()* 9.99 + 0.01).toFixed(2)
}

// https://k6.io/docs/using-k6/options#list-of-options
export let options = {
    // A number specifying the number of VUs to run concurrently
    vus: 250,
    // A list of objects that specify the target number of VUs to ramp up or down;
    // shortcut option for a single scenario with a ramping VUs executor
    stages: [
        { duration: "40s", target: 300 },
        // { duration: "60s", target: 2500 },
        // { duration: "120s", target: 8000 },
        // { duration: "120s", target: 2000 },
        // { duration: "60s", target: 2500 },
        // { duration: "40s", target: 1000 },
        // { duration: "40s", target: 200 },
    ],

    // A boolean specifying whether to throw errors on failed HTTP requests
    throw: true,
};

const traderURL = __ENV.TRADER_URL ? __ENV.TRADER_URL : 'http://localhost:4000'

export default function() {
    let req_data = {
        "item_id": randItemID(),
        "order_type": randOrderType(),
        "price": randPrice(),
    };
    
    // Send a JSON encoded POST request
    let body = JSON.stringify(req_data);
    // console.log(body)
    let headers = { 
        "Content-Type": "application/json",
        "X-Identity-Id": "1",
    }
    let res = http.post(`${traderURL}/apis/v1/orders`, body, { headers: headers});

    // Verify response
    check(res, {
        "status is 200": (r) => r.status === 200,
    });
}
