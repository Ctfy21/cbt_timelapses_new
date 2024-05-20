import { defineStore } from "pinia";
import { Order } from "@/class/order";
import {socket} from "@/socket";

export const useOrderStore = defineStore("orders",{
    state: () => ({
        orders: [],
        folders: [],
        error: ""
    }),

    actions: {
        bindEvents(){
            socket.onopen = () => {
                console.log("Successfully connected!");
            }
            socket.onmessage = (event) => {
                let json = JSON.parse(event.data)
                console.log(json)
                for (const [key, value] of Object.entries(json)){
                    switch (key){
                        case "Orders": {
                            const ordersArray = [...Order.generateOrder(value)]
                            console.log(ordersArray.sort((a,b) => b.id - a.id ))
                            this.orders = ordersArray.sort((a,b) => b.id - a.id )
                            break
                        }
                        case "Order": {
                            const order = Order.deserialize(value);
                            console.log(order);
                            const existingOrderIndex = this.orders.findIndex((temp) => {
                                return temp.id === order.id;
                            });
                            if (existingOrderIndex != -1) {
                                this.orders[existingOrderIndex] = order
                            }
                            break
                        }
                        case "Folders": {
                            if (value !== null){
                                this.folders = value;
                            }
                            break
                        }
                        default:
                            console.log("Error during parse JSON")
                    }
                }
            }
        },
        addNewOrder(room, camera, startDate, endDate){
            this.error = ''
            let new_order = Order.createOrder(room, camera, startDate, endDate)
            if(new_order !== undefined){
                if(socket.readyState === WebSocket.OPEN){
                    socket.send(Order.serialize(new_order))
                }
            }
            else{
                this.error = "Неправильно указана дата и время или неправильно указан период"
            }
        },
    }
})