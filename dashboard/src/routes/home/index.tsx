import { FunctionalComponent, h } from 'preact';
import { useEffect, useRef, useState } from 'preact/hooks';
import feedPB from '../../protobuf/feed_pb';
import style from './style.css';


const Home: FunctionalComponent = () => {
    const [address, setAddress] = useState('ws://localhost:32001/ws')
    const [connectionState, setConnectionState] = useState<number | undefined>(undefined)
    const [velocity, setVelocity] = useState<number[] | undefined>(undefined)
    const [position, setPosition] = useState<number[] | undefined>(undefined)
    const [fsu, setFSU] = useState<number[] | undefined>(undefined)
    const websocketRef = useRef<WebSocket | undefined>()
    const dataReceived = useRef(0)
    const [dataRate, setDataRate] = useState(0)

    useEffect(() => {
        let prevDataReceived = 0
        let prevTime = 0
        const interval = setInterval(() => {
            const timeNow = +new Date()
            const timeDelta = (timeNow - prevTime) / 1000
            setDataRate((dataReceived.current - prevDataReceived) / timeDelta)
            prevTime = timeNow
            prevDataReceived = dataReceived.current
        }, 1000)
        return () => {
            clearInterval(interval)
        }
    }, [])

    const onButtonClick = () => {
        if (websocketRef.current) {
            websocketRef.current.close()
            websocketRef.current = undefined;
        }

        websocketRef.current = new WebSocket(address)
        websocketRef.current.binaryType = 'arraybuffer'
        setConnectionState(websocketRef.current.readyState)
        websocketRef.current.onopen = () => {
            if (websocketRef.current) {
                setConnectionState(websocketRef.current.readyState)
            }
        }
        websocketRef.current.onerror = () => {
            if (websocketRef.current) {
                setConnectionState(websocketRef.current.readyState)
            }
        }
        websocketRef.current.onclose = () => {
            if (websocketRef.current) {
                setConnectionState(websocketRef.current.readyState)
            }
            websocketRef.current = undefined
        }
        websocketRef.current.onmessage = (e) => {
            dataReceived.current += e.data.byteLength;
            const pmove = feedPB.PMove.deserializeBinary(e.data);
            setVelocity(pmove.getVelocityList())
            setPosition(pmove.getPositionList())
            setFSU(pmove.getFsuList())
        }
    }

    const onAddressInput: h.JSX.GenericEventHandler<HTMLInputElement> = (e) => {
        if (e.target instanceof HTMLInputElement) {
            setAddress(e.target.value)
        }
    }

    const connectionStateText = (() => {
        if (typeof connectionState !== 'number') {
            return 'Idle'
        }
        switch (connectionState) {
            case WebSocket.CONNECTING:
                return 'Connecting';
            case WebSocket.OPEN:
                return 'Open';
            case WebSocket.CLOSING:
                return 'Closing';
            case WebSocket.CLOSED:
                return 'Closed';
            default:
                return 'Unknown';
        }
    })()

    return (
        <div class={style.home}>
            <h1>Home</h1>
            <p>This is the Home component.</p>
            <label>
                Address: <input type="text" value={address} onInput={onAddressInput} />
            </label>
            <div>Connection state: {connectionStateText}</div>
            <div>Downlink data rate: {(dataRate / 1000).toFixed(3)} KB/s</div>
            <button onClick={onButtonClick}>Connect</button>

            {velocity ? <div>Velocity: {velocity[0].toFixed(6)} {velocity[1].toFixed(6)} {velocity[2].toFixed(6)}</div> : null}
            {position ? <div>Position: {position[0].toFixed(6)} {position[1].toFixed(6)} {position[2].toFixed(6)}</div> : null}
            {fsu ? <div>FSU: {fsu[0].toFixed(6)} {fsu[1].toFixed(6)} {fsu[2].toFixed(6)}</div> : null}
        </div>
    );
};

export default Home;
