export class Order {

    static status_ok = 200
    static status_waiting = 300
    static status_error = 400


    constructor(room, camera, startDate, endDate, status) {
        this.room = room
        this.camera = camera;
        this.startDate = new Date(startDate);
        this.endDate = new Date(endDate);
        this.status = status;
    }


    static serialize(order){
        let toJSONOrder = {
            Order: {
                'Room': order.room,
                'Camera': order.camera,
                'StartDate': order.startDate.toISOString().slice(0,10) + "_" + "00-00-00",
                'EndDate': order.endDate.toISOString().slice(0,10) + "_" + "00-00-00",
                'Status': order.status,
            }
        }
        return JSON.stringify(toJSONOrder)
    }

    static deserialize(json){
        const order = new Order();
        order.id = json["Id"]
        order.room = json["Room"]
        order.camera = json["Camera"];
        order.startDate = new Date(json["StartDate"].split("_")[0]);
        order.endDate = new Date(json["EndDate"].split("_")[0]);
        order.status = json["Status"];
        order.downloaderValue = ''
        return order;
    }

    static createOrder(room, camera, startDate, endDate) {
        const new_order = new this(room, camera, startDate, endDate, Order.status_waiting)
        if(endDate <= startDate){
            return undefined
        }
        console.log(this.serialize(new_order))
        return new_order
    }

    static *generateOrder(arrayJsonOrder){
        for(const value of arrayJsonOrder){
            let order = JSON.parse(value)
            yield Order.deserialize(order["Order"]);
        }
    }
}



