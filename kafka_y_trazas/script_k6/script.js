import { Writer } from 'k6/x/kafka'
import { sleep } from 'k6'
import encoding from 'k6/encoding'
 
export const options = { vus: 1, duration: '1h' }
 
const writer = new Writer({
  brokers: [__ENV.KAFKA_BROKER],
  topic: "pruebas-k6",
  autoCreateTopic: true,
  requiredAcks: 1,
})
 
export default function () {
 
  const payload = JSON.stringify({
    ts: Date.now(),
    msg: "hola"
  })
 
  const msg = {
    value: encoding.b64encode(payload)
  }
 
  writer.produce({
    messages: [msg]
  })
 
  console.log("Mensaje enviado OK")
 
  sleep(3)
}