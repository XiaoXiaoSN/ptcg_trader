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

export let options = {
    vus: 25,
    stages: [
        { duration: "40s", target: 1000 },
        // { duration: "60s", target: 2500 },
        // { duration: "120s", target: 8000 },
        // { duration: "120s", target: 2000 },
        // { duration: "60s", target: 2500 },
        // { duration: "40s", target: 1000 },
        // { duration: "40s", target: 200 },
    ]
};

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
    let res = http.post("http://localhost:4000/apis/v1/orders", body, { headers: headers});

    // Verify response
    check(res, {
        "status is 200": (r) => r.status === 200,
        "status is 409": (r) => r.status === 409,
        "status is 500": (r) => r.status === 500,
    });
}
