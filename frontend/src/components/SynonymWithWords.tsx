import React from 'react'

interface WordWithSynonymsProps {
  synonym: string
  words: string[]
}

export const SynonymWithWords: React.FC<WordWithSynonymsProps> = ({
  words,
  synonym,
}) => {
  return (
    <div className="text-center flex gap-4 flex-col">
      <h1 className="text-4xl text-primary">{synonym}</h1>
      <div className="flex justify-center">
        <p className="text-white">{'['}&nbsp;</p>
        {words.map((word, index) => (
          <span key={index} className="text-md text-green-300">
            {index > 0 && <span className="text-white">,&nbsp;</span>}
            {word}
          </span>
        ))}
        <p className="text-white">&nbsp;{']'}</p>
      </div>
    </div>
  )
}
