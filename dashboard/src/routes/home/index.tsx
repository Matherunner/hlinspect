import { FunctionalComponent, h } from 'preact';
import { useRef, useState } from 'preact/hooks';
import feedPB from '../../protobuf/feed_pb';
import style from './style.css';


const Home: FunctionalComponent = () => {
    const [address, setAddress] = useState('ws://localhost:32001/ws')
    const [connectionState, setConnectionState] = useState<number | undefined>(undefined)
    const [velocity, setVelocity] = useState<number[] | undefined>(undefined)
    const websocketRef = useRef<WebSocket | undefined>()

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
            const pmove = feedPB.PMove.deserializeBinary(e.data);
            setVelocity(pmove.getVelocityList())
        }
    }

    const onAddressInput: h.JSX.GenericEventHandler<HTMLInputElement> = (e) => {
        if (e.target instanceof HTMLInputElement) {
            setAddress(e.target.value)
        }
    }

    const connectionStateText = (() => {
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
            <button onClick={onButtonClick}>Connect</button>

            {velocity ? <div>Velocity: {velocity[0]} {velocity[1]} {velocity[2]}</div> : null}
        </div>
    );
};

export default Home;
